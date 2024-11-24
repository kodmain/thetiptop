package services

import (
	"github.com/gofiber/fiber/v2"
	transfert "github.com/kodmain/thetiptop/api/internal/application/transfert/user"
	"github.com/kodmain/thetiptop/api/internal/application/validator"
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
	if err := dtoClient.Check(data.Validator{
		"id": {validator.Required, validator.ID},
	}); err != nil {
		return err.Code(), err
	}

	// Attempt to delete the client using the service
	if err := service.DeleteClient(dtoClient); err != nil {
		return err.Code(), err
	}

	// Return 204 if the deletion is successful with no content
	return fiber.StatusNoContent, nil
}

func GetClient(service services.UserServiceInterface, dtoClient *transfert.Client) (int, any) {
	if err := dtoClient.Check(data.Validator{
		"id": {validator.Required, validator.ID},
	}); err != nil {
		return err.Code(), err
	}

	client, err := service.GetClient(dtoClient)
	if err != nil {
		return err.Code(), err
	}

	return fiber.StatusOK, client
}

func UpdateClient(service services.UserServiceInterface, clientDTO *transfert.Client) (int, any) {
	if err := clientDTO.Check(data.Validator{
		"id":         {validator.Required, validator.ID},
		"newsletter": {validator.IsBool},
	}); err != nil {
		return err.Code(), err
	}

	client, err := service.UpdateClient(clientDTO)
	if err != nil {
		return err.Code(), err
	}

	return fiber.StatusOK, client
}

func RegisterClient(service services.UserServiceInterface, credentialDTO *transfert.Credential, clientDTO *transfert.Client) (int, any) {
	if err := credentialDTO.Check(data.Validator{
		"email":    {validator.Required, validator.Email},
		"password": {validator.Required, validator.Password},
	}); err != nil {
		return err.Code(), err
	}

	if err := clientDTO.Check(data.Validator{
		"newsletter": {validator.Required, validator.IsBool},
		"cgu":        {validator.Required, validator.IsBool, validator.IsTrue},
	}); err != nil {
		return err.Code(), err
	}

	credential, err := service.RegisterClient(credentialDTO, clientDTO)
	if err != nil {
		return err.Code(), err
	}

	return fiber.StatusCreated, credential
}
