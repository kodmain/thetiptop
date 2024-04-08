package workflow

import (
	"github.com/gofiber/fiber/v2"
	"github.com/kodmain/thetiptop/api/internal/application/dto"
	"github.com/kodmain/thetiptop/api/internal/domain/services"
)

var Client = &ClientWorkflow{}

type ClientWorkflow struct{}

func (c *ClientWorkflow) SignUp(email, password string) (int, fiber.Map) {
	dto, err := dto.NewClient(email, password)
	if err != nil {
		return fiber.StatusBadRequest, fiber.Map{"error": err.Error()}
	}

	if err := services.Client().SignUp(dto); err != nil {
		return fiber.StatusBadRequest, fiber.Map{"error": err.Error()}
	}

	return fiber.StatusCreated, nil
}
