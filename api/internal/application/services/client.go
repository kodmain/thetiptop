package services

import (
	"github.com/gofiber/fiber/v2"
	"github.com/kodmain/thetiptop/api/internal/application/transfert"
	"github.com/kodmain/thetiptop/api/internal/application/validator"
	"github.com/kodmain/thetiptop/api/internal/domain/client/errors"
	"github.com/kodmain/thetiptop/api/internal/domain/client/services"
	"github.com/kodmain/thetiptop/api/internal/infrastructure/data"
)

func SignValidation(service services.ClientServiceInterface, dtoValidation *transfert.Validation, dtoCredential *transfert.Credential) (int, any) {
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

	validation, err := service.SignValidation(dtoValidation, dtoCredential)
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

func ValidationRecover(service services.ClientServiceInterface, dtoCredential *transfert.Credential, dtoValidation *transfert.Validation) (int, any) {
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

func UpdateClient(service services.ClientServiceInterface, clientDTO *transfert.Client) (int, any) {
	err := clientDTO.Check(data.Validator{
		"id":         {validator.Required, validator.ID},
		"newsletter": {validator.IsBool},
	})

	if err != nil {
		return fiber.StatusBadRequest, err.Error()
	}

	err = service.UpdateClient(clientDTO)
	if err != nil {
		return fiber.StatusInternalServerError, err.Error()
	}

	return fiber.StatusNoContent, nil
}
