package services

import (
	"github.com/google/uuid"
	"github.com/kodmain/thetiptop/api/internal/domain/entities"
	"github.com/kodmain/thetiptop/api/internal/domain/repositories"
)

type ClientService struct {
	repo repositories.ClientRepository
}

func NewClientService(repo repositories.ClientRepository) *ClientService {
	return &ClientService{
		repo: repo,
	}
}

func (s *ClientService) GetClient(uuid uuid.UUID) (*entities.Client, error) {
	return s.repo.GetClient(uuid)
}

func (s *ClientService) GetClients(filter map[string]interface{}) ([]*entities.Client, error) {
	return s.repo.GetClients(filter)
}

func (s *ClientService) CreateClient(client *entities.Client) (*entities.Client, error) {
	return s.repo.CreateClient(client)
}

func (s *ClientService) UpdateClient(client *entities.Client) (*entities.Client, error) {
	return s.repo.UpdateClient(client)
}

func (s *ClientService) DeleteClient(uuid uuid.UUID) error {
	return s.repo.DeleteClient(uuid)
}
