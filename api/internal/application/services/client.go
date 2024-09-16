package services

import (
	"github.com/gofiber/fiber/v2"
	"github.com/kodmain/thetiptop/api/internal/application/transfert"
	"github.com/kodmain/thetiptop/api/internal/application/validator"
	"github.com/kodmain/thetiptop/api/internal/domain/user/errors"
	"github.com/kodmain/thetiptop/api/internal/domain/user/services"
	"github.com/kodmain/thetiptop/api/internal/infrastructure/data"
)

// DeleteClient deletes a client by ID
// This function checks the client's ID, validates it, and proceeds to delete the client from the system
//
// Parameters:
// - service: services.UserServiceInterface The service responsible for client management
// - dtoClient: *transfert.Client The DTO that contains the client's data to be deleted
//
// Returns:
// - int: The HTTP status code indicating success or failure of the operation
// - any: The response object, which can be an error message in case of failure, or nil for successful deletion
func DeleteClient(service services.UserServiceInterface, dtoClient *transfert.Client) (int, any) {
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

func GetClient(service services.UserServiceInterface, dtoClient *transfert.Client) (int, any) {
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

func UpdateClient(service services.UserServiceInterface, clientDTO *transfert.Client) (int, any) {
	err := clientDTO.Check(data.Validator{
		"id":         {validator.Required, validator.ID},
		"newsletter": {validator.IsBool},
	})

	if err != nil {
		return fiber.StatusBadRequest, err.Error()
	}

	client, err := service.UpdateClient(clientDTO)
	if err != nil {
		if err.Error() == errors.ErrClientNotFound {
			return fiber.StatusNotFound, err.Error()
		}

		return fiber.StatusInternalServerError, err.Error()
	}

	return fiber.StatusOK, client
}

func RegisterClient(service services.UserServiceInterface, credentialDTO *transfert.Credential, clientDTO *transfert.Client) (int, any) {
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

	credential, err := service.RegisterClient(credentialDTO, clientDTO)
	if err != nil {
		if err.Error() == errors.ErrCredentialAlreadyExists {
			return fiber.StatusConflict, err.Error()
		}

		return fiber.StatusInternalServerError, err.Error()
	}

	return fiber.StatusCreated, credential
}
