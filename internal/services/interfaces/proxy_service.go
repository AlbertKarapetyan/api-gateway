package interfaces

import "net/http"

type ProxyService interface {
	InitServers()
	ReverseProxy(w http.ResponseWriter, r *http.Request)
	HealthCheck()
}
