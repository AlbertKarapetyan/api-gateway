package services

import (
	"api-gateway/internal/config"
	"api-gateway/internal/models"
	"api-gateway/internal/services/interfaces"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"sync"
	"time"
)

var (
	servers []*models.Server
	mu      sync.Mutex
)

type proxyService struct {
	lb interfaces.LoadBalancer
	bs interfaces.BackEndService
}

func NewProxyService(lb interfaces.LoadBalancer, bs interfaces.BackEndService) interfaces.ProxyService {
	return &proxyService{lb: lb, bs: bs}
}

func (p *proxyService) InitServers() {
	mu.Lock()
	defer mu.Unlock()

	servers = []*models.Server{}
	for _, urlStr := range config.CFG.Servers {
		parsedURL, err := url.Parse(urlStr)
		if err != nil {
			log.Printf("❌ Invalid server URL: %s\n", urlStr)
			continue
		}

		// Perform health check
		healthCheckURL := parsedURL.String() + "/health"
		isAlive := checkServerHealth(healthCheckURL)

		server := &models.Server{URL: parsedURL, IsAlive: isAlive}
		servers = append(servers, server)
	}
}

func (p *proxyService) ReverseProxy(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/favicon.ico" {
		http.NotFound(w, r)
		return
	}

	server := p.lb.GetNextServer(servers)
	if server == nil {
		http.Error(w, "No servers available", http.StatusServiceUnavailable)
		return
	}

	proxy := httputil.NewSingleHostReverseProxy(server.URL)
	p.lb.IncrementConnections(server)
	r.Host = server.URL.Host

	proxy.ModifyResponse = func(resp *http.Response) error {
		p.lb.DecrementConnections(server)
		return nil
	}

	proxy.ServeHTTP(w, r)
}

func (p *proxyService) HealthCheck() {
	for {
		for _, backend := range servers {
			healthCheckURL := backend.URL.String() + "/health"
			isAlive := checkServerHealth(healthCheckURL)
			p.bs.SetServer(backend)
			p.bs.SetAlive(isAlive)
			if isAlive {
				log.Printf("✅ Backend %s is UP\n", backend.URL)
			} else {
				log.Printf("❌ Backend %s is DOWN\n", backend.URL)
			}
		}

		time.Sleep(time.Duration(config.CFG.HealthCheckInterval) * time.Second) // checking every HealthCheckInterval seconds
	}
}

func checkServerHealth(url string) bool {
	rs, err := http.Get(url)
	if err != nil || rs.StatusCode != http.StatusOK {
		if rs != nil {
			rs.Body.Close()
		}
		return false
	}
	rs.Body.Close()
	return true
}
