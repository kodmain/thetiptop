package repositories

import (
	"fmt"

	"github.com/kodmain/thetiptop/api/internal/application/transfert"
	"github.com/kodmain/thetiptop/api/internal/domain/entities"
	"github.com/kodmain/thetiptop/api/internal/domain/errors"
	"github.com/kodmain/thetiptop/api/internal/infrastructure/observability/logger"
	"github.com/kodmain/thetiptop/api/internal/infrastructure/providers/database"
	"gorm.io/gorm"
)

type ClientRepository struct {
	database *gorm.DB
}

func NewClientRepository(name string) *ClientRepository {
	database := database.Get(name)
	database.AutoMigrate(&entities.Client{})

	return &ClientRepository{database}
}

func (r *ClientRepository) Create(obj *transfert.Client) (*entities.Client, error) {
	client, err := entities.CreateClient(obj)
	if err != nil {
		return nil, err
	}

	result := r.database.Create(client)

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

	result := r.database.Where(obj).First(client)

	logger.Error(result.Error)

	if result.Error != nil {
		return nil, result.Error
	}

	return client, nil
}
