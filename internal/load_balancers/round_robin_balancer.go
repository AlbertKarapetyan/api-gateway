package loadbalancers

import (
	"api-gateway/internal/models"
	"api-gateway/internal/services/interfaces"
	"sync"
	"sync/atomic"
)

var index uint64

type roundRobinBalancer struct {
	connections sync.Map
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

// IncrementConnections implements interfaces.LoadBalancer.
func (b *roundRobinBalancer) IncrementConnections(server *models.Server) {
	val, _ := b.connections.LoadOrStore(server.URL, int32(0))
	b.connections.Store(server.URL, val.(int32)+1)
}

// DecrementConnections implements interfaces.LoadBalancer.
func (b *roundRobinBalancer) DecrementConnections(server *models.Server) {
	val, _ := b.connections.LoadOrStore(server.URL, int32(0))
	if val.(int32) > 0 {
		b.connections.Store(server.URL, val.(int32)-1)
	}
}
