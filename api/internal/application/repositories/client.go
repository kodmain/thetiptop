package repositories

import (
	"github.com/google/uuid"
	"github.com/kodmain/thetiptop/api/internal/architecture/persistence"
	"github.com/kodmain/thetiptop/api/internal/domain/entities"
	"gorm.io/gorm"
)

type ClientRepository struct {
	database *gorm.DB
}

func NewClientRepository() *ClientRepository {
	return &ClientRepository{
		database: persistence.Get("file"),
	}
}

// CreateClient implements repositories.ClientRepository.
func (r ClientRepository) CreateClient(client *entities.Client) (*entities.Client, error) {
	panic("unimplemented")
}

// DeleteClient implements repositories.ClientRepository.
func (r ClientRepository) DeleteClient(id uuid.UUID) error {
	panic("unimplemented")
}

// GetClient implements repositories.ClientRepository.
func (r ClientRepository) GetClient(id uuid.UUID) (*entities.Client, error) {
	panic("unimplemented")
}

// GetClients implements repositories.ClientRepository.
func (r ClientRepository) GetClients(filter map[string]any) ([]*entities.Client, error) {
	panic("unimplemented")
}

// UpdateClient implements repositories.ClientRepository.
func (r ClientRepository) UpdateClient(client *entities.Client) (*entities.Client, error) {
	panic("unimplemented")
}
