package services

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/kodmain/thetiptop/api/internal/application/transfert"
	"github.com/kodmain/thetiptop/api/internal/application/validator"
	"github.com/kodmain/thetiptop/api/internal/domain/user/errors"
	"github.com/kodmain/thetiptop/api/internal/domain/user/services"
	"github.com/kodmain/thetiptop/api/internal/infrastructure/data"
	serializer "github.com/kodmain/thetiptop/api/internal/infrastructure/serializers/jwt"
)

func UserAuth(service services.UserServiceInterface, credentialDTO *transfert.Credential) (int, any) {
	err := credentialDTO.Check(data.Validator{
		"email":    {validator.Required, validator.Email},
		"password": {validator.Required, validator.Password},
	})

	if err != nil {
		return fiber.StatusBadRequest, err.Error()
	}

	userID, err := service.UserAuth(credentialDTO)
	if err != nil {
		return fiber.StatusBadRequest, err.Error()
	}

	accessToken, refreshToken, err := serializer.FromID(*userID)
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

func CredentialUpdate(service services.UserServiceInterface, validationDTO *transfert.Validation, credentialDTO *transfert.Credential) (int, any) {
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

func MailValidation(service services.UserServiceInterface, dtoValidation *transfert.Validation, dtoCredential *transfert.Credential) (int, any) {
	err := dtoValidation.Check(data.Validator{
		"token": {validator.Required, validator.Luhn},
	})

	if err != nil {
		return fiber.StatusBadRequest, err.Error()
	}

	err = dtoCredential.Check(data.Validator{
		"email": {validator.Required, validator.Email},
	})

	if err != nil {
		return fiber.StatusBadRequest, err.Error()
	}

	validation, err := service.MailValidation(dtoValidation, dtoCredential)
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

func ValidationRecover(service services.UserServiceInterface, dtoCredential *transfert.Credential, dtoValidation *transfert.Validation) (int, any) {
	err := dtoCredential.Check(data.Validator{
		"email": {validator.Required, validator.Email},
	})

	if err != nil {
		return fiber.StatusBadRequest, err.Error()
	}

	err = dtoValidation.Check(data.Validator{
		"type": {validator.Required},
	})

	if err != nil {
		return fiber.StatusBadRequest, err.Error()
	}

	if err = service.ValidationRecover(dtoValidation, dtoCredential); err != nil {
		if err.Error() == errors.ErrUserNotFound {
			return fiber.StatusNotFound, err.Error()
		}

		return fiber.StatusBadRequest, err.Error()
	}

	return fiber.StatusNoContent, nil
}
