package services

import (
	"github.com/gofiber/fiber/v2"
	"github.com/kodmain/thetiptop/api/internal/domain/code/services"
)

func ListErrors(service services.CodeServiceInterface) (int, any) {
	errors, err := service.ListErrors()
	if err != nil {
		return err.Code(), err
	}

	return fiber.StatusOK, errors
}
