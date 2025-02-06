package interfaces

import "api-gateway/internal/models"

type BackEndService interface {
	SetServer(server *models.Server)
	SetAlive(isAlive bool)
	IsAlive() bool
}
