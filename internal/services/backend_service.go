package services

import (
	"api-gateway/internal/models"
	"api-gateway/internal/services/interfaces"
	"log"
	"sync"
)

type backEndService struct {
	server *models.Server

	mu sync.RWMutex
}

func NewBackEndService() interfaces.BackEndService {
	return &backEndService{server: nil, mu: sync.RWMutex{}}
}

func (b *backEndService) SetServer(server *models.Server) {
	b.server = server
}

func (b *backEndService) SetAlive(isAlive bool) {
	if b.server != nil {
		b.mu.Lock()
		b.server.IsAlive = isAlive
		b.mu.Unlock()
	} else {
		log.Fatal("The instance of server is nil, before this method call the SetServer()")
	}
}

func (b *backEndService) IsAlive() bool {
	if b.server != nil {
		b.mu.RLock()
		defer b.mu.RUnlock()
		return b.server.IsAlive
	} else {
		log.Fatal("The instance of server is nil, before this method call the SetServer()")
	}
	return false
}
