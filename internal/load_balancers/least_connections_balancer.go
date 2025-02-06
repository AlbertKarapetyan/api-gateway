package loadbalancers

import (
	"api-gateway/internal/models"
	"api-gateway/internal/services/interfaces"
	"sync"
)

type leastConnectionsBalancer struct {
	connections sync.Map
}

func NewLeastConnectionsBalancer() interfaces.LoadBalancer {
	return &leastConnectionsBalancer{}
}

// GetNextServer implements interfaces.LoadBalancer.
func (l *leastConnectionsBalancer) GetNextServer(servers []*models.Server) *models.Server {
	var selected *models.Server
	minConnections := int32(1<<31 - 1) // maximum int value

	for _, server := range servers {
		if !server.IsAlive {
			continue
		}

		cnt, ok := l.connections.Load(server.URL)
		connections := int32(0)
		if ok {
			connections = cnt.(int32)
		}

		if connections < minConnections {
			selected = server
			minConnections = connections
		}
	}

	return selected
}

// IncrementConnections implements interfaces.LoadBalancer.
func (l *leastConnectionsBalancer) IncrementConnections(server *models.Server) {
	val, _ := l.connections.LoadOrStore(server.URL, int32(0))
	l.connections.Store(server.URL, val.(int32)+1)
}

// DecrementConnections implements interfaces.LoadBalancer.
func (l *leastConnectionsBalancer) DecrementConnections(server *models.Server) {
	val, _ := l.connections.LoadOrStore(server.URL, int32(0))
	if val.(int32) > 0 {
		l.connections.Store(server.URL, val.(int32)-1)
	}
}
