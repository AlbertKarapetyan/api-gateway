package interfaces

import "api-gateway/internal/models"

type LoadBalancer interface {
	GetNextServer(servers []*models.Server) *models.Server
	IncrementConnections(server *models.Server)
	DecrementConnections(server *models.Server)
}
