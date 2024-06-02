package services

import (
	"github.com/gofiber/fiber/v2"
	"github.com/kodmain/thetiptop/api/internal/application/transfert"
	"github.com/kodmain/thetiptop/api/internal/domain/errors"
	"github.com/kodmain/thetiptop/api/internal/domain/services"
	"github.com/kodmain/thetiptop/api/internal/infrastructure/data"
	"github.com/kodmain/thetiptop/api/internal/infrastructure/providers/database"
	"github.com/kodmain/thetiptop/api/internal/infrastructure/providers/mail"
	"github.com/kodmain/thetiptop/api/internal/infrastructure/repositories"

	serializer "github.com/kodmain/thetiptop/api/internal/infrastructure/serializers/jwt"
)

func SignUp(email, password string) (int, fiber.Map) {
	obj, err := transfert.NewClient(data.Object{
		"email":    email,
		"password": password,
	})

	if err != nil {
		return fiber.StatusBadRequest, fiber.Map{"error": err.Error()}
	}

	clientService := services.Client(
		repositories.NewClientRepository(database.Get()),
		mail.Get(),
	)

	if err != nil {
		return fiber.StatusInternalServerError, fiber.Map{"error": err.Error()}
	}

	client, err := clientService.SignUp(obj)
	if err != nil {
		if err.Error() == errors.ErrClientAlreadyExists {
			return fiber.StatusConflict, fiber.Map{"error": err.Error()}
		}

		return fiber.StatusInternalServerError, fiber.Map{"error": err.Error()}
	}

	return fiber.StatusCreated, fiber.Map{"client": client}
}

func SignIn(email, password string) (int, fiber.Map) {
	obj, err := transfert.NewClient(data.Object{
		"email":    email,
		"password": password,
	})

	if err != nil {
		return fiber.StatusBadRequest, fiber.Map{"error": err.Error()}
	}

	clientService := services.Client(
		repositories.NewClientRepository(database.Get()),
		mail.Get(),
	)

	if err != nil {
		return fiber.StatusInternalServerError, fiber.Map{"error": err.Error()}
	}

	client, err := clientService.SignIn(obj)
	if err != nil {
		return fiber.StatusBadRequest, fiber.Map{"error": err.Error()}
	}

	token, err := serializer.FromID(client.ID)
	if err != nil {
		return fiber.StatusInternalServerError, fiber.Map{"error": err.Error()}
	}

	return fiber.StatusOK, fiber.Map{"jwt": token}
}

func SignRenew(refresh *serializer.Token) (int, fiber.Map) {
	if refresh == nil {
		return fiber.StatusBadRequest, fiber.Map{"error": "Invalid token"}
	}

	if refresh.Type != serializer.REFRESH {
		return fiber.StatusBadRequest, fiber.Map{"error": "Invalid token type"}
	}

	if refresh.HasExpired() {
		return fiber.StatusUnauthorized, fiber.Map{"error": "Refresh token has expired"}
	}

	refreshToken, err := serializer.FromID(refresh.ID)
	if err != nil {
		return fiber.StatusInternalServerError, fiber.Map{"error": err.Error()}
	}

	return fiber.StatusOK, fiber.Map{"jwt": refreshToken}
}
