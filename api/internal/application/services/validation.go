package services

import (
	"github.com/gofiber/fiber/v2"
	"github.com/kodmain/thetiptop/api/internal/application/transfert"
	"github.com/kodmain/thetiptop/api/internal/application/validator"
	"github.com/kodmain/thetiptop/api/internal/domain/user/errors"
	"github.com/kodmain/thetiptop/api/internal/domain/user/services"
	"github.com/kodmain/thetiptop/api/internal/infrastructure/data"
)

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
		if err.Error() == errors.ErrClientNotFound {
			return fiber.StatusNotFound, err.Error()
		}

		return fiber.StatusBadRequest, err.Error()
	}

	return fiber.StatusNoContent, nil
}
