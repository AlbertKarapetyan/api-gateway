package main

import (
	"api-gateway/internal/app"
	"api-gateway/internal/config"
	"log"
	"net/http"
)

func main() {
	config.LoadConfigs()

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
