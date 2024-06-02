package repositories

import (
	"fmt"
	"sync"

	"github.com/kodmain/thetiptop/api/internal/application/transfert"
	"github.com/kodmain/thetiptop/api/internal/domain/entities"
	"github.com/kodmain/thetiptop/api/internal/domain/errors"
	"github.com/kodmain/thetiptop/api/internal/infrastructure/providers/database"
)

var once sync.Once

type ClientRepository struct {
	store *database.Database
}

func NewClientRepository(store *database.Database) *ClientRepository {
	once.Do(func() {
		store.Engine.AutoMigrate(entities.Client{})
	})

	return &ClientRepository{store}
}

func (r *ClientRepository) Create(obj *transfert.Client) (*entities.Client, error) {
	client := entities.CreateClient(obj)
	result := r.store.Engine.Create(client)

	if result.Error != nil {
		if result.Error.Error() == "UNIQUE constraint failed: clients.email" {
			return nil, fmt.Errorf(errors.ErrClientAlreadyExists)
		}

		return nil, result.Error
	}

	return client, nil
}

func (r *ClientRepository) Read(obj *transfert.Client) (*entities.Client, error) {
	client := entities.CreateClient(obj)
	result := r.store.Engine.Where(obj).First(client)

	if result.Error != nil {
		return nil, result.Error
	}

	return client, nil
}
