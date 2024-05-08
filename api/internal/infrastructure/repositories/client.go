package repositories

import (
	"fmt"

	"github.com/kodmain/thetiptop/api/internal/application/transfert"
	"github.com/kodmain/thetiptop/api/internal/domain/entities"
	"github.com/kodmain/thetiptop/api/internal/domain/errors"
	"github.com/kodmain/thetiptop/api/internal/infrastructure/providers/database"
)

type ClientRepository struct {
	service database.ServiceInterface
}

func NewClientRepository(name string) *ClientRepository {
	service := database.Get(name)
	service.AutoMigrate(&entities.Client{})

	return &ClientRepository{service}
}

func (r *ClientRepository) Create(obj *transfert.Client) (*entities.Client, error) {
	client, err := entities.CreateClient(obj)
	if err != nil {
		return nil, err
	}

	result := r.service.Create(client)

	if result.Error != nil {
		if result.Error.Error() == "UNIQUE constraint failed: clients.email" {
			return nil, fmt.Errorf(errors.ErrClientAlreadyExists)
		}

		return nil, result.Error
	}

	return client, nil
}

func (r *ClientRepository) Read(obj *transfert.Client) (*entities.Client, error) {
	client, err := entities.CreateClient(obj)
	if err != nil {
		return nil, err
	}

	result := r.service.Where(obj).First(client)

	if result.Error != nil {
		return nil, result.Error
	}

	return client, nil
}
