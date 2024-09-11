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

func UserRegister(service services.ClientServiceInterface, credentialDTO *transfert.Credential, clientDTO *transfert.Client) (int, any) {
	err := credentialDTO.Check(data.Validator{
		"email":    {validator.Required, validator.Email},
		"password": {validator.Required, validator.Password},
	})

	if err != nil {
		return fiber.StatusBadRequest, err.Error()
	}

	err = clientDTO.Check(data.Validator{
		"newsletter": {validator.Required, validator.IsBool},
		"cgu":        {validator.Required, validator.IsBool, validator.IsTrue},
	})

	if err != nil {
		return fiber.StatusBadRequest, err.Error()
	}

	credential, err := service.UserRegister(credentialDTO, clientDTO)
	if err != nil {
		if err.Error() == errors.ErrCredentialAlreadyExists {
			return fiber.StatusConflict, err.Error()
		}

		return fiber.StatusInternalServerError, err.Error()
	}

	return fiber.StatusCreated, credential
}

func UserAuth(service services.ClientServiceInterface, credentialDTO *transfert.Credential) (int, any) {
	err := credentialDTO.Check(data.Validator{
		"email":    {validator.Required, validator.Email},
		"password": {validator.Required, validator.Password},
	})

	if err != nil {
		return fiber.StatusBadRequest, err.Error()
	}

	client, err := service.UserAuth(credentialDTO)
	if err != nil {
		return fiber.StatusBadRequest, err.Error()
	}

	accessToken, refreshToken, err := serializer.FromID(*client.Credential.ClientID)
	if err != nil {
		return fiber.StatusInternalServerError, err.Error()
	}

	return fiber.StatusOK, fiber.Map{
		"access_token":  accessToken,
		"refresh_token": refreshToken,
	}
}

func UserAuthRenew(refresh *serializer.Token) (int, any) {
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

func CredentialUpdate(service services.ClientServiceInterface, validationDTO *transfert.Validation, credentialDTO *transfert.Credential) (int, any) {
	err := validationDTO.Check(data.Validator{
		"token": {validator.Required, validator.Luhn},
	})

	if err != nil {
		return fiber.StatusBadRequest, err.Error()
	}

	err = credentialDTO.Check(data.Validator{
		"email":    {validator.Required, validator.Email},
		"password": {validator.Required, validator.Password},
	})

	if err != nil {
		return fiber.StatusBadRequest, err.Error()
	}

	validation, err := service.PasswordValidation(validationDTO, &transfert.Credential{
		Email: credentialDTO.Email,
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

	err = service.PasswordUpdate(credentialDTO)
	if err != nil {
		return fiber.StatusInternalServerError, err.Error()
	}

	return fiber.StatusOK, validation
}
