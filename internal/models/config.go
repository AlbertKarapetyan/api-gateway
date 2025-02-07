package models

type Config struct {
	BalancerType        string              `json:"balancer_type"`
	Servers             map[string][]string `json:"servers"`
	HealthCheckInterval int                 `json:"health_check_interval"`
	Routes              map[string]string   `json:"routes"`
}
