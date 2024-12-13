package store

import (
	"github.com/gofiber/fiber/v2"
	"github.com/kodmain/thetiptop/api/internal/domain/store/services"
)

func ListStores(service services.StoreServiceInterface) (int, any) {
	credential, err := service.ListStores()
	if err != nil {
		return err.Code(), err
	}

	return fiber.StatusOK, credential
}
