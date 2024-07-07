package services

import (
	"github.com/gofiber/fiber/v2"
	"github.com/kodmain/thetiptop/api/internal/application/transfert"
	"github.com/kodmain/thetiptop/api/internal/application/validator"
	"github.com/kodmain/thetiptop/api/internal/domain/client/errors"
	"github.com/kodmain/thetiptop/api/internal/domain/client/repositories"
	"github.com/kodmain/thetiptop/api/internal/domain/client/services"
	"github.com/kodmain/thetiptop/api/internal/infrastructure/data"
	"github.com/kodmain/thetiptop/api/internal/infrastructure/providers/database"
	"github.com/kodmain/thetiptop/api/internal/infrastructure/providers/mail"
	serializer "github.com/kodmain/thetiptop/api/internal/infrastructure/serializers/jwt"
)

func SignUp(service services.ClientServiceInterface, email, password string) (int, fiber.Map) {
	obj, err := transfert.NewClient(data.Object{
		"email":    &email,
		"password": &password,
	}, data.Validator{
		"email":    {validator.Required, validator.Email},
		"password": {validator.Required, validator.Password},
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

func SignIn(service services.ClientServiceInterface, email, password string) (int, fiber.Map) {
	obj, err := transfert.NewClient(data.Object{
		"email":    &email,
		"password": &password,
	}, data.Validator{
		"email":    {validator.Required, validator.Email},
		"password": {validator.Required, validator.Password},
	})

	if err != nil {
		return fiber.StatusBadRequest, fiber.Map{"error": err.Error()}
	}

	client, err := service.SignIn(obj)
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

func SignValidation(email, token string) (int, fiber.Map) {
	obj, err := transfert.NewValidation(data.Object{
		"token": &token,
		"email": &email,
	}, data.Validator{
		"token": {validator.Required, validator.Luhn},
		"email": {validator.Required, validator.ID},
	})

	if err != nil {
		return fiber.StatusBadRequest, fiber.Map{"error": err.Error()}
	}

	clientService := services.Client(
		repositories.NewClientRepository(database.Get()),
		mail.Get(),
	)

	validation, err := clientService.SignValidation(obj)
	if err != nil {
		status := fiber.StatusInternalServerError
		switch err.Error() {
		case errors.ErrValidationNotFound:
			status = fiber.StatusNotFound
		case errors.ErrValidationAlreadyValidated:
			status = fiber.StatusConflict
		case errors.ErrValidationExpired:
			status = fiber.StatusGone
		}

		return status, fiber.Map{"error": err.Error()}
	}

	return fiber.StatusOK, fiber.Map{"validation": validation}

}

func PasswordRecover(service services.ClientServiceInterface, email string) (int, fiber.Map) {
	obj, err := transfert.NewClient(data.Object{
		"email": &email,
	}, data.Validator{
		"email": {validator.Required, validator.Email},
	})

	if err != nil {
		return fiber.StatusBadRequest, fiber.Map{"error": err.Error()}
	}

	if err = service.PasswordRecover(obj); err != nil {
		return fiber.StatusBadRequest, fiber.Map{"error": err.Error()}
	}

	return fiber.StatusNoContent, nil
}

func PasswordValidation(clientId, token, password string) (int, fiber.Map) {
	dtoClient, err := transfert.NewClient(data.Object{
		"password": &password,
		"id":       &clientId,
	}, data.Validator{
		"password": {validator.Required, validator.Password},
		"id":       {validator.Required, validator.ID},
	})

	if err != nil {
		return fiber.StatusBadRequest, fiber.Map{"error": err.Error()}
	}

	dtoValidation, err := transfert.NewValidation(data.Object{
		"token":    &token,
		"clientId": &clientId,
	}, data.Validator{
		"token":    {validator.Required, validator.Luhn},
		"clientId": {validator.Required, validator.ID},
	})

	if err != nil {
		return fiber.StatusBadRequest, fiber.Map{"error": err.Error()}
	}

	clientService := services.Client(
		repositories.NewClientRepository(database.Get()),
		mail.Get(),
	)

	validation, err := clientService.PasswordValidation(dtoValidation)
	if err != nil {
		status := fiber.StatusInternalServerError
		switch err.Error() {
		case errors.ErrValidationNotFound:
			status = fiber.StatusNotFound
		case errors.ErrValidationAlreadyValidated:
			status = fiber.StatusConflict
		case errors.ErrValidationExpired:
			status = fiber.StatusGone
		}

		return status, fiber.Map{"error": err.Error()}
	}

	err = clientService.PasswordUpdate(dtoClient)
	if err != nil {
		return fiber.StatusInternalServerError, fiber.Map{"error": err.Error()}
	}

	return fiber.StatusOK, fiber.Map{"validation": validation}
}
