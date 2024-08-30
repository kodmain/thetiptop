package services

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/kodmain/thetiptop/api/internal/application/transfert"
	"github.com/kodmain/thetiptop/api/internal/application/validator"
	"github.com/kodmain/thetiptop/api/internal/domain/client/errors"
	"github.com/kodmain/thetiptop/api/internal/domain/client/services"
	"github.com/kodmain/thetiptop/api/internal/infrastructure/data"
	serializer "github.com/kodmain/thetiptop/api/internal/infrastructure/serializers/jwt"
)

func SignUp(service services.ClientServiceInterface, clientDTO *transfert.Client) (int, any) {
	err := clientDTO.Check(data.Validator{
		"email":    {validator.Required, validator.Email},
		"password": {validator.Required, validator.Password},
	})

	if err != nil {
		return fiber.StatusBadRequest, err
	}

	client, err := service.SignUp(clientDTO)
	if err != nil {
		if err.Error() == errors.ErrClientAlreadyExists {
			return fiber.StatusConflict, err
		}

		return fiber.StatusInternalServerError, err
	}

	return fiber.StatusCreated, client
}

func SignIn(service services.ClientServiceInterface, clientDTO *transfert.Client) (int, any) {
	err := clientDTO.Check(data.Validator{
		"email":    {validator.Required, validator.Email},
		"password": {validator.Required, validator.Password},
	})

	if err != nil {
		return fiber.StatusBadRequest, err
	}

	client, err := service.SignIn(clientDTO)
	if err != nil {
		return fiber.StatusBadRequest, err
	}

	token, err := serializer.FromID(client.ID)
	if err != nil {
		return fiber.StatusInternalServerError, err
	}

	return fiber.StatusOK, token
}

func SignRenew(refresh *serializer.Token) (int, any) {
	if refresh == nil {
		return fiber.StatusBadRequest, fmt.Errorf("invalid token")
	}

	if refresh.Type != serializer.REFRESH {
		return fiber.StatusBadRequest, fmt.Errorf("invalid token type")
	}

	if refresh.HasExpired() {
		return fiber.StatusUnauthorized, fmt.Errorf("refresh token has expired")
	}

	refreshToken, err := serializer.FromID(refresh.ID)
	if err != nil {
		return fiber.StatusInternalServerError, err
	}

	return fiber.StatusOK, refreshToken
}

func SignValidation(service services.ClientServiceInterface, validationDTO *transfert.Validation, clientDTO *transfert.Client) (int, any) {
	err := validationDTO.Check(data.Validator{
		"token": {validator.Required, validator.Luhn},
	})

	if err != nil {
		return fiber.StatusBadRequest, err
	}

	err = clientDTO.Check(data.Validator{
		"email": {validator.Required, validator.Email},
	})

	if err != nil {
		return fiber.StatusBadRequest, err
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

		return status, err
	}

	return fiber.StatusOK, validation

}

func PasswordRecover(service services.ClientServiceInterface, clientDTO *transfert.Client) (int, any) {
	err := clientDTO.Check(data.Validator{
		"email": {validator.Required, validator.Email},
	})

	if err != nil {
		return fiber.StatusBadRequest, err
	}

	if err = service.PasswordRecover(clientDTO); err != nil {
		if err.Error() == errors.ErrClientNotFound {
			return fiber.StatusNotFound, err
		}

		return fiber.StatusBadRequest, err
	}

	return fiber.StatusNoContent, nil
}

func PasswordUpdate(service services.ClientServiceInterface, validationDTO *transfert.Validation, clientDTO *transfert.Client) (int, any) {
	err := validationDTO.Check(data.Validator{
		"token": {validator.Required, validator.Luhn},
	})

	if err != nil {
		return fiber.StatusBadRequest, err
	}

	err = clientDTO.Check(data.Validator{
		"email":    {validator.Required, validator.Email},
		"password": {validator.Required, validator.Password},
	})

	if err != nil {
		return fiber.StatusBadRequest, err
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

		return status, err
	}

	err = service.PasswordUpdate(clientDTO)
	if err != nil {
		return fiber.StatusInternalServerError, err
	}

	return fiber.StatusOK, validation
}
