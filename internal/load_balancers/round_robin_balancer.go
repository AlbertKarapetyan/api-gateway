package loadbalancers

import (
	"api-gateway/internal/models"
	"api-gateway/internal/services/interfaces"
	"sync/atomic"
)

var index uint64

type roundRobinBalancer struct {
	BaseLoadBalancer
}

func NewRoundRobinBalancer() interfaces.LoadBalancer {
	return &roundRobinBalancer{}
}

// GetNextServer implements interfaces.LoadBalancer.
func (b *roundRobinBalancer) GetNextServer(servers []*models.Server) *models.Server {
	for i := 0; i < len(servers); i++ {
		next := atomic.AddUint64(&index, 1) % uint64(len(servers))
		server := servers[next]
		if server.IsAlive {
			return server
		}
	}
	return nil
}
