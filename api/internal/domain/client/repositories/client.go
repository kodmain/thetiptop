package repositories

import (
	"fmt"
	"sync"

	"github.com/kodmain/thetiptop/api/internal/application/transfert"
	"github.com/kodmain/thetiptop/api/internal/domain/client/entities"
	"github.com/kodmain/thetiptop/api/internal/domain/client/errors"
	"github.com/kodmain/thetiptop/api/internal/infrastructure/providers/database"
)

var client sync.Once

type ClientRepository struct {
	store *database.Database
}

type ClientRepositoryInterface interface {
	// client
	CreateClient(obj *transfert.Client) (*entities.Client, error)
	ReadClient(obj *transfert.Client) (*entities.Client, error)
	UpdateClient(entity *entities.Client) error
	DeleteClient(obj *transfert.Client) error

	// validation
	CreateValidation(obj *transfert.Validation) (*entities.Validation, error)
	ReadValidation(obj *transfert.Validation) (*entities.Validation, error)
	UpdateValidation(entity *entities.Validation) error
	DeleteValidation(obj *transfert.Validation) error
}

func NewClientRepository(store *database.Database) *ClientRepository {
	client.Do(func() {
		store.Engine.AutoMigrate(entities.Client{})
		store.Engine.AutoMigrate(entities.Validation{})
	})

	return &ClientRepository{store}
}

func (r *ClientRepository) CreateClient(obj *transfert.Client) (*entities.Client, error) {
	client := entities.CreateClient(obj)
	client.Validations = append(client.Validations, &entities.Validation{
		Type: entities.Mail,
	})

	result := r.store.Engine.Create(client)
	if result.Error != nil {
		if result.Error.Error() == "UNIQUE constraint failed: clients.email" {
			return nil, fmt.Errorf(errors.ErrClientAlreadyExists)
		}

		return nil, result.Error
	}

	return client, nil
}

func (r *ClientRepository) ReadClient(obj *transfert.Client) (*entities.Client, error) {
	client := entities.CreateClient(obj)
	result := r.store.Engine.Where(obj).First(client)

	if result.Error != nil {
		return nil, result.Error
	}

	return client, nil
}

func (r *ClientRepository) UpdateClient(entity *entities.Client) error {
	return r.store.Engine.Save(entity).Error
}

func (r *ClientRepository) DeleteClient(obj *transfert.Client) error {
	client := entities.CreateClient(obj)
	result := r.store.Engine.Delete(client)

	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (r *ClientRepository) CreateValidation(obj *transfert.Validation) (*entities.Validation, error) {
	validation := entities.CreateValidation(obj)
	result := r.store.Engine.Create(validation)

	if result.Error != nil {
		return nil, result.Error
	}

	return validation, nil
}

func (r *ClientRepository) ReadValidation(obj *transfert.Validation) (*entities.Validation, error) {
	validation := entities.CreateValidation(obj)
	result := r.store.Engine.Where(obj).First(validation)

	if result.Error != nil {
		return nil, result.Error
	}

	return validation, nil
}

func (r *ClientRepository) UpdateValidation(entity *entities.Validation) error {
	return r.store.Engine.Save(entity).Error
}

func (r *ClientRepository) DeleteValidation(obj *transfert.Validation) error {
	validation := entities.CreateValidation(obj)
	result := r.store.Engine.Delete(validation)

	if result.Error != nil {
		return result.Error
	}

	return nil
}
