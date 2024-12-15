package store

import (
	"github.com/gofiber/fiber/v2"
	transfert "github.com/kodmain/thetiptop/api/internal/application/transfert/crm"
	"github.com/kodmain/thetiptop/api/internal/application/validator"
	"github.com/kodmain/thetiptop/api/internal/domain/store/services"
	"github.com/kodmain/thetiptop/api/internal/infrastructure/data"
)

func ListStores(service services.StoreServiceInterface) (int, any) {
	credential, err := service.ListStores()
	if err != nil {
		return err.Code(), err
	}

	return fiber.StatusOK, credential
}

func GetStoreByID(service services.StoreServiceInterface, dtoStore *transfert.Store) (int, any) {
	if err := dtoStore.Check(data.Validator{
		"id": {validator.Required, validator.ID},
	}); err != nil {
		return err.Code(), err
	}

	store, err := service.GetStoreByID(dtoStore)
	if err != nil {
		return err.Code(), err
	}

	return fiber.StatusOK, store
}
