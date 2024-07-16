package services

import (
	"github.com/gofiber/fiber/v2"
	"github.com/kodmain/thetiptop/api/internal/application/transfert"
	"github.com/kodmain/thetiptop/api/internal/application/validator"
	"github.com/kodmain/thetiptop/api/internal/domain/client/errors"
	"github.com/kodmain/thetiptop/api/internal/domain/client/services"
	"github.com/kodmain/thetiptop/api/internal/infrastructure/data"
	serializer "github.com/kodmain/thetiptop/api/internal/infrastructure/serializers/jwt"
)

func SignUp(service services.ClientServiceInterface, clientDTO *transfert.Client) (int, fiber.Map) {
	err := clientDTO.Check(data.Validator{
		"email":    {validator.Required, validator.Email},
		"password": {validator.Required, validator.Password},
	})

	if err != nil {
		return fiber.StatusBadRequest, fiber.Map{"error": err.Error()}
	}

	client, err := service.SignUp(clientDTO)
	if err != nil {
		if err.Error() == errors.ErrClientAlreadyExists {
			return fiber.StatusConflict, fiber.Map{"error": err.Error()}
		}

		return fiber.StatusInternalServerError, fiber.Map{"error": err.Error()}
	}

	return fiber.StatusCreated, fiber.Map{"client": client}
}

func SignIn(service services.ClientServiceInterface, clientDTO *transfert.Client) (int, fiber.Map) {
	err := clientDTO.Check(data.Validator{
		"email":    {validator.Required, validator.Email},
		"password": {validator.Required, validator.Password},
	})

	if err != nil {
		return fiber.StatusBadRequest, fiber.Map{"error": err.Error()}
	}

	client, err := service.SignIn(clientDTO)
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

func SignValidation(service services.ClientServiceInterface, validationDTO *transfert.Validation, clientDTO *transfert.Client) (int, fiber.Map) {
	err := validationDTO.Check(data.Validator{
		"token": {validator.Required, validator.Luhn},
	})

	if err != nil {
		return fiber.StatusBadRequest, fiber.Map{"error": err.Error()}
	}

	err = clientDTO.Check(data.Validator{
		"email": {validator.Required, validator.Email},
	})

	if err != nil {
		return fiber.StatusBadRequest, fiber.Map{"error": err.Error()}
	}

	validation, err := service.SignValidation(validationDTO, clientDTO)
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

func PasswordRecover(service services.ClientServiceInterface, clientDTO *transfert.Client) (int, fiber.Map) {
	err := clientDTO.Check(data.Validator{
		"email": {validator.Required, validator.Email},
	})

	if err != nil {
		return fiber.StatusBadRequest, fiber.Map{"error": err.Error()}
	}

	if err = service.PasswordRecover(clientDTO); err != nil {
		if err.Error() == errors.ErrClientNotFound {
			return fiber.StatusNotFound, fiber.Map{"error": err.Error()}
		}

		return fiber.StatusBadRequest, fiber.Map{"error": err.Error()}
	}

	return fiber.StatusNoContent, nil
}

func PasswordUpdate(service services.ClientServiceInterface, validationDTO *transfert.Validation, clientDTO *transfert.Client) (int, fiber.Map) {
	err := validationDTO.Check(data.Validator{
		"token": {validator.Required, validator.Luhn},
	})

	if err != nil {
		return fiber.StatusBadRequest, fiber.Map{"error": err.Error()}
	}

	err = clientDTO.Check(data.Validator{
		"email":    {validator.Required, validator.Email},
		"password": {validator.Required, validator.Password},
	})

	if err != nil {
		return fiber.StatusBadRequest, fiber.Map{"error": err.Error()}
	}

	validation, err := service.PasswordValidation(validationDTO, &transfert.Client{
		Email: clientDTO.Email,
	})

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

	err = service.PasswordUpdate(clientDTO)
	if err != nil {
		return fiber.StatusInternalServerError, fiber.Map{"error": err.Error()}
	}

	return fiber.StatusOK, fiber.Map{"validation": validation}
}
