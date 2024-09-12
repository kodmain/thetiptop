package services

import (
	"github.com/gofiber/fiber/v2"
	"github.com/kodmain/thetiptop/api/internal/application/transfert"
	"github.com/kodmain/thetiptop/api/internal/application/validator"
	"github.com/kodmain/thetiptop/api/internal/domain/client/errors"
	"github.com/kodmain/thetiptop/api/internal/domain/client/services"
	"github.com/kodmain/thetiptop/api/internal/infrastructure/data"
)

// DeleteClient deletes a client by ID
// This function checks the client's ID, validates it, and proceeds to delete the client from the system
//
// Parameters:
// - service: services.ClientServiceInterface The service responsible for client management
// - dtoClient: *transfert.Client The DTO that contains the client's data to be deleted
//
// Returns:
// - int: The HTTP status code indicating success or failure of the operation
// - any: The response object, which can be an error message in case of failure, or nil for successful deletion
func DeleteClient(service services.ClientServiceInterface, dtoClient *transfert.Client) (int, any) {
	// Validation of the client ID
	err := dtoClient.Check(data.Validator{
		"id": {validator.Required, validator.ID},
	})

	// Return 400 if validation fails
	if err != nil {
		return fiber.StatusBadRequest, err.Error()
	}

	// Attempt to delete the client using the service
	err = service.DeleteClient(dtoClient)
	if err != nil {
		// Handle specific error cases like client not found
		if err.Error() == errors.ErrClientNotFound {
			return fiber.StatusNotFound, err.Error()
		}

		// Return 500 for any internal server errors
		return fiber.StatusInternalServerError, err.Error()
	}

	// Return 204 if the deletion is successful with no content
	return fiber.StatusNoContent, nil
}

func GetClient(service services.ClientServiceInterface, dtoClient *transfert.Client) (int, any) {
	err := dtoClient.Check(data.Validator{
		"id": {validator.Required, validator.ID},
	})

	if err != nil {
		return fiber.StatusBadRequest, err.Error()
	}

	client, err := service.GetClient(dtoClient)
	if err != nil {
		if err.Error() == errors.ErrClientNotFound {
			return fiber.StatusNotFound, err.Error()
		}

		return fiber.StatusInternalServerError, err.Error()
	}

	return fiber.StatusOK, client
}

func MailValidation(service services.ClientServiceInterface, dtoValidation *transfert.Validation, dtoCredential *transfert.Credential) (int, any) {
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
