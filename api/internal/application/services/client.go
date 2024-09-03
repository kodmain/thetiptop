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
		"email":      {validator.Required, validator.Email},
		"password":   {validator.Required, validator.Password},
		"newsletter": {validator.Required, validator.IsBool},
		"cgu":        {validator.Required, validator.IsBool, validator.IsTrue},
	})

	if err != nil {
		return fiber.StatusBadRequest, err.Error()
	}

	client, err := service.SignUp(clientDTO)
	if err != nil {
		if err.Error() == errors.ErrClientAlreadyExists {
			return fiber.StatusConflict, err.Error()
		}

		return fiber.StatusInternalServerError, err.Error()
	}

	return fiber.StatusCreated, client
}

func SignIn(service services.ClientServiceInterface, clientDTO *transfert.Client) (int, any) {
	err := clientDTO.Check(data.Validator{
		"email":    {validator.Required, validator.Email},
		"password": {validator.Required, validator.Password},
	})

	if err != nil {
		return fiber.StatusBadRequest, err.Error()
	}

	client, err := service.SignIn(clientDTO)
	if err != nil {
		return fiber.StatusBadRequest, err.Error()
	}

	accessToken, refreshToken, err := serializer.FromID(client.ID)
	if err != nil {
		return fiber.StatusInternalServerError, err.Error()
	}

	return fiber.StatusOK, fiber.Map{
		"access_token":  accessToken,
		"refresh_token": refreshToken,
	}
}

func SignRenew(refresh *serializer.Token) (int, any) {
	if refresh == nil {
		return fiber.StatusBadRequest, fmt.Errorf("invalid token").Error()
	}

	if refresh.Type != serializer.REFRESH {
		return fiber.StatusBadRequest, fmt.Errorf("invalid token type").Error()
	}

	if refresh.HasExpired() {
		return fiber.StatusUnauthorized, fmt.Errorf("refresh token has expired").Error()
	}

	accessToken, refreshToken, err := serializer.FromID(refresh.ID)
	if err != nil {
		return fiber.StatusInternalServerError, err.Error()
	}

	return fiber.StatusOK, fiber.Map{
		"access_token":  accessToken,
		"refresh_token": refreshToken,
	}
}

func SignValidation(service services.ClientServiceInterface, validationDTO *transfert.Validation, clientDTO *transfert.Client) (int, any) {
	err := validationDTO.Check(data.Validator{
		"token": {validator.Required, validator.Luhn},
	})

	if err != nil {
		return fiber.StatusBadRequest, err.Error()
	}

	err = clientDTO.Check(data.Validator{
		"email": {validator.Required, validator.Email},
	})

	if err != nil {
		return fiber.StatusBadRequest, err.Error()
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

		return status, err.Error()
	}

	return fiber.StatusOK, validation

}

func PasswordRecover(service services.ClientServiceInterface, clientDTO *transfert.Client) (int, any) {
	err := clientDTO.Check(data.Validator{
		"email": {validator.Required, validator.Email},
	})

	if err != nil {
		return fiber.StatusBadRequest, err.Error()
	}

	if err = service.PasswordRecover(clientDTO); err != nil {
		if err.Error() == errors.ErrClientNotFound {
			return fiber.StatusNotFound, err.Error()
		}

		return fiber.StatusBadRequest, err.Error()
	}

	return fiber.StatusNoContent, nil
}

func ValidationRecover(service services.ClientServiceInterface, clientDTO *transfert.Client, validationDTO *transfert.Validation) (int, any) {
	err := clientDTO.Check(data.Validator{
		"email": {validator.Required, validator.Email},
	})

	if err != nil {
		return fiber.StatusBadRequest, err.Error()
	}

	err = validationDTO.Check(data.Validator{
		"type": {validator.Required},
	})

	if err != nil {
		return fiber.StatusBadRequest, err.Error()
	}

	if err = service.ValidationRecover(validationDTO, clientDTO); err != nil {
		if err.Error() == errors.ErrClientNotFound {
			return fiber.StatusNotFound, err.Error()
		}

		return fiber.StatusBadRequest, err.Error()
	}

	return fiber.StatusNoContent, nil
}

func PasswordUpdate(service services.ClientServiceInterface, validationDTO *transfert.Validation, clientDTO *transfert.Client) (int, any) {
	err := validationDTO.Check(data.Validator{
		"token": {validator.Required, validator.Luhn},
	})

	if err != nil {
		return fiber.StatusBadRequest, err.Error()
	}

	err = clientDTO.Check(data.Validator{
		"email":    {validator.Required, validator.Email},
		"password": {validator.Required, validator.Password},
	})

	if err != nil {
		return fiber.StatusBadRequest, err.Error()
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

		return status, err.Error()
	}

	err = service.PasswordUpdate(clientDTO)
	if err != nil {
		return fiber.StatusInternalServerError, err.Error()
	}

	return fiber.StatusOK, validation
}
