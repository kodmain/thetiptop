package services

import (
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/gofiber/fiber/v2"
	"github.com/kodmain/thetiptop/api/internal/application/transfert"
	"github.com/kodmain/thetiptop/api/internal/application/validator"
	"github.com/kodmain/thetiptop/api/internal/domain/client/errors"
	"github.com/kodmain/thetiptop/api/internal/domain/client/services"
	"github.com/kodmain/thetiptop/api/internal/infrastructure/data"
)

func PasswordRecover(service services.ClientServiceInterface, email string) (int, fiber.Map) {
	obj, err := transfert.NewClient(data.Object{
		"email": aws.String(email),
	}, data.Validator{
		"email": {validator.Required, validator.Email},
	})

	if err != nil {
		return fiber.StatusBadRequest, fiber.Map{"error": err.Error()}
	}

	client, err := service.SignUp(obj)
	if err != nil {
		if err.Error() == errors.ErrClientAlreadyExists {
			return fiber.StatusConflict, fiber.Map{"error": err.Error()}
		}

		return fiber.StatusInternalServerError, fiber.Map{"error": err.Error()}
	}

	return fiber.StatusCreated, fiber.Map{"client": client}
}
