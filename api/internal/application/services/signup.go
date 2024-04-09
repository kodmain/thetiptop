package services

import (
	"github.com/gofiber/fiber/v2"
	"github.com/kodmain/thetiptop/api/internal/application/dto"
	"github.com/kodmain/thetiptop/api/internal/domain/errors"
	"github.com/kodmain/thetiptop/api/internal/domain/services"
)

func SignUp(email, password string) (int, fiber.Map) {
	dto, err := dto.NewClient(email, password)
	if err != nil {
		return fiber.StatusBadRequest, fiber.Map{"error": err.Error()}
	}

	if err := services.Client().SignUp(dto); err != nil {
		if err.Error() == errors.ErrClientAlreadyExists {
			return fiber.StatusConflict, fiber.Map{"error": err.Error()}
		}

		return fiber.StatusBadRequest, fiber.Map{"error": err.Error()}
	}

	return fiber.StatusCreated, nil
}
