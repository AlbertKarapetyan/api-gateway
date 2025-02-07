package main

import (
	"api-gateway/internal/app"
	"api-gateway/internal/config"
	"log"
	"net/http"
)

func main() {
	config.LoadConfigs()

	// Print config for debugging
	log.Printf("Balancer Type: %s", config.CFG.BalancerType)
	log.Printf("Health Check Interval: %d", config.CFG.HealthCheckInterval)

	log.Println("Servers:")
	for service, servers := range config.CFG.Servers {
		log.Printf("  %s: %v", service, servers)
	}

	log.Println("Routes:")
	for route, path := range config.CFG.Routes {
		log.Printf("  %s -> %s", route, path)
	}

	configChangeChan := make(chan struct{})
	go config.CheckingConfigs(configChangeChan)

	proxyService, err := app.InitProxyService()
	if err != nil {
		log.Fatal("❌ Error init ProxyService:", err)
	}

	proxyService.InitServers()
	go proxyService.HealthCheck()

	http.HandleFunc("/", proxyService.ReverseProxy)

	log.Println("🚀 Load Balancer started on :8080")
	go func() {
		log.Fatal(http.ListenAndServe(":8080", nil))
	}()

	// Слушаем канал изменений конфигурации
	for range configChangeChan {
		log.Println("🔄 Config.json change detected, reinit servers...")

		config.LoadConfigs()
		proxyService.InitServers()
	}
}
