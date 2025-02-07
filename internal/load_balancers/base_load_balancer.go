package loadbalancers

import (
	"api-gateway/internal/models"
	"sync"
)

// BaseLoadBalancer contains common logic for managing server connections
type BaseLoadBalancer struct {
	connections sync.Map
}

// IncrementConnections increments the connection count for a given server
func (b *BaseLoadBalancer) IncrementConnections(server *models.Server) {
	val, _ := b.connections.LoadOrStore(server.URL, int32(0))
	b.connections.Store(server.URL, val.(int32)+1)
}

// DecrementConnections decrements the connection count for a given server
func (b *BaseLoadBalancer) DecrementConnections(server *models.Server) {
	val, _ := b.connections.LoadOrStore(server.URL, int32(0))
	if val.(int32) > 0 {
		b.connections.Store(server.URL, val.(int32)-1)
	}
}
