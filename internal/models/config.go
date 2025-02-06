package models

type Config struct {
	BalancerType        string   `json:"balancer_type"`
	Servers             []string `json:"servers"`
	HealthCheckInterval int      `json:"health_check_interval"`
}
