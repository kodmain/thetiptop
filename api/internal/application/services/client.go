package services

import (
	"github.com/gofiber/fiber/v2"
	"github.com/kodmain/thetiptop/api/internal/application/transfert"
	"github.com/kodmain/thetiptop/api/internal/domain/errors"
	"github.com/kodmain/thetiptop/api/internal/domain/services"
	"github.com/kodmain/thetiptop/api/internal/infrastructure/data"
	"github.com/kodmain/thetiptop/api/internal/infrastructure/observability/logger"

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

	if err := services.Client().SignUp(obj); err != nil {
		if err.Error() == errors.ErrClientAlreadyExists {
			return fiber.StatusConflict, fiber.Map{"error": err.Error()}
		}

		return fiber.StatusBadRequest, fiber.Map{"error": err.Error()}
	}

	return fiber.StatusCreated, nil
}

func SignRenew(token *serializer.Token) (int, fiber.Map) {
	logger.Infof("Access: %v %v %v", token.Exp, token.TZ, token.Offset)

	refresh := token.HasRefresh()

	if refresh == nil {
		return fiber.StatusBadRequest, fiber.Map{"error": "No refresh token"}
	}

	logger.Infof("Refresh: %v %v %v", refresh.Exp, refresh.TZ, refresh.Offset)

	if refresh.HasExpired() {
		return fiber.StatusUnauthorized, fiber.Map{"error": "Refresh token has expired"}
	}

	refreshToken, err := serializer.FromID(refresh.ID)
	if err != nil {
		return fiber.StatusInternalServerError, fiber.Map{"error": err.Error()}
	}

	return fiber.StatusOK, fiber.Map{"jwt": refreshToken}
}

func SignIn(email, password string) (int, fiber.Map) {
	obj, err := transfert.NewClient(data.Object{
		"email":    email,
		"password": password,
	})

	if err != nil {
		return fiber.StatusBadRequest, fiber.Map{"error": err.Error()}
	}

	client, err := services.Client().SignIn(obj)
	if err != nil {
		return fiber.StatusBadRequest, fiber.Map{"error": err.Error()}
	}

	token, err := serializer.FromID(client.ID)
	if err != nil {
		return fiber.StatusInternalServerError, fiber.Map{"error": err.Error()}
	}

	return fiber.StatusOK, fiber.Map{"jwt": token}
}
