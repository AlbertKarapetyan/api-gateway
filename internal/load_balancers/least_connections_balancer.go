package loadbalancers

import (
	"api-gateway/internal/models"
	"api-gateway/internal/services/interfaces"
)

type leastConnectionsBalancer struct {
	BaseLoadBalancer
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
