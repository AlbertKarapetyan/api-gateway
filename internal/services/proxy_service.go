package services

import (
	"api-gateway/internal/config"
	"api-gateway/internal/models"
	"api-gateway/internal/services/interfaces"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
	"sync"
	"time"
)

var (
	servers map[string][]*models.Server
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

	servers = make(map[string][]*models.Server) // Initialize the map

	for serviceType, serverURLs := range config.CFG.Servers {
		var serviceServers []*models.Server

		for _, urlStr := range serverURLs {
			parsedURL, err := url.Parse(urlStr)
			if err != nil {
				log.Printf("❌ Invalid server URL for service %s: %s\n", serviceType, urlStr)
				continue
			}

			// Perform health check
			healthCheckURL := parsedURL.String() + "/health"
			isAlive := checkServerHealth(healthCheckURL)

			server := &models.Server{URL: parsedURL, IsAlive: isAlive}
			serviceServers = append(serviceServers, server)
		}
		servers[serviceType] = serviceServers
	}
}

func (p *proxyService) HealthCheck() {
	for {
		for serviceType, serviceServers := range servers {
			for _, backend := range serviceServers {
				healthCheckURL := backend.URL.String() + "/health"
				isAlive := checkServerHealth(healthCheckURL)
				p.bs.SetServer(backend)
				p.bs.SetAlive(isAlive)
				if isAlive {
					log.Printf("✅ Backend %s (%s) is UP\n", backend.URL, serviceType)
				} else {
					log.Printf("❌ Backend %s (%s) is DOWN\n", backend.URL, serviceType)
				}
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

func (p *proxyService) ReverseProxy(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/favicon.ico" {
		http.NotFound(w, r)
		return
	}

	// Determine the service type based on the route
	serviceType := getServiceType(r.URL.Path)
	if serviceType == "" {
		http.Error(w, "Invalid route", http.StatusNotFound)
		return
	}

	// Get the list of servers for the service type
	serviceServers, ok := servers[serviceType]
	if !ok || len(serviceServers) == 0 {
		http.Error(w, "No servers configured for this service", http.StatusServiceUnavailable)
		return
	}

	// Construct the target URL using the backend path from the Routes config
	backendPath := config.CFG.Routes[r.URL.Path]
	if backendPath == "" {
		http.Error(w, "Invalid backend path", http.StatusInternalServerError)
		return
	}

	// Get the next available server using the load balancer
	server := p.lb.GetNextServer(serviceServers)
	if server == nil {
		http.Error(w, "No servers available", http.StatusServiceUnavailable)
		return
	}

	targetURL := server.URL.String() + backendPath
	parsedURL, err := url.Parse(targetURL)
	if err != nil {
		http.Error(w, "Failed to parse target URL", http.StatusInternalServerError)
		return
	}

	proxy := httputil.NewSingleHostReverseProxy(parsedURL)
	p.lb.IncrementConnections(server)
	r.Host = server.URL.Host

	proxy.ModifyResponse = func(resp *http.Response) error {
		p.lb.DecrementConnections(server)
		return nil
	}

	proxy.ServeHTTP(w, r)
}

func getServiceType(path string) string {
	for route := range config.CFG.Routes {
		if strings.HasPrefix(path, route) {
			// Extract the service type from the route (e.g., "/user/signin" -> "user")
			parts := strings.Split(route, "/")
			if len(parts) > 0 {
				return parts[1]
			}
		}
	}
	return ""
}
