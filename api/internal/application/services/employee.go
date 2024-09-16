package services

import (
	"github.com/gofiber/fiber/v2"
	"github.com/kodmain/thetiptop/api/internal/application/transfert"
	"github.com/kodmain/thetiptop/api/internal/application/validator"
	"github.com/kodmain/thetiptop/api/internal/domain/user/errors"
	"github.com/kodmain/thetiptop/api/internal/domain/user/services"
	"github.com/kodmain/thetiptop/api/internal/infrastructure/data"
)

// DeleteEmployee deletes a employee by ID
// This function checks the employee's ID, validates it, and proceeds to delete the employee from the system
//
// Parameters:
// - service: services.UserServiceInterface The service responsible for employee management
// - dtoEmployee: *transfert.Employee The DTO that contains the employee's data to be deleted
//
// Returns:
// - int: The HTTP status code indicating success or failure of the operation
// - any: The response object, which can be an error message in case of failure, or nil for successful deletion
func DeleteEmployee(service services.UserServiceInterface, dtoEmployee *transfert.Employee) (int, any) {
	// Validation of the employee ID
	err := dtoEmployee.Check(data.Validator{
		"id": {validator.Required, validator.ID},
	})

	// Return 400 if validation fails
	if err != nil {
		return fiber.StatusBadRequest, err.Error()
	}

	// Attempt to delete the employee using the service
	err = service.DeleteEmployee(dtoEmployee)
	if err != nil {
		// Handle specific error cases like employee not found
		if err.Error() == errors.ErrEmployeeNotFound {
			return fiber.StatusNotFound, err.Error()
		}

		// Return 500 for any internal server errors
		return fiber.StatusInternalServerError, err.Error()
	}

	// Return 204 if the deletion is successful with no content
	return fiber.StatusNoContent, nil
}

func GetEmployee(service services.UserServiceInterface, dtoEmployee *transfert.Employee) (int, any) {
	err := dtoEmployee.Check(data.Validator{
		"id": {validator.Required, validator.ID},
	})

	if err != nil {
		return fiber.StatusBadRequest, err.Error()
	}

	employee, err := service.GetEmployee(dtoEmployee)
	if err != nil {
		if err.Error() == errors.ErrEmployeeNotFound {
			return fiber.StatusNotFound, err.Error()
		}

		return fiber.StatusInternalServerError, err.Error()
	}

	return fiber.StatusOK, employee
}

func UpdateEmployee(service services.UserServiceInterface, employeeDTO *transfert.Employee) (int, any) {
	err := employeeDTO.Check(data.Validator{
		"id": {validator.Required, validator.ID},
	})

	if err != nil {
		return fiber.StatusBadRequest, err.Error()
	}

	employee, err := service.UpdateEmployee(employeeDTO)
	if err != nil {
		return fiber.StatusInternalServerError, err.Error()
	}

	return fiber.StatusNoContent, employee
}

func RegisterEmployee(service services.UserServiceInterface, credentialDTO *transfert.Credential, employeeDTO *transfert.Employee) (int, any) {
	err := credentialDTO.Check(data.Validator{
		"email":    {validator.Required, validator.Email},
		"password": {validator.Required, validator.Password},
	})

	if err != nil {
		return fiber.StatusBadRequest, err.Error()
	}

	err = employeeDTO.Check(data.Validator{})

	if err != nil {
		return fiber.StatusBadRequest, err.Error()
	}

	credential, err := service.RegisterEmployee(credentialDTO, employeeDTO)
	if err != nil {
		if err.Error() == errors.ErrCredentialAlreadyExists {
			return fiber.StatusConflict, err.Error()
		}

		return fiber.StatusInternalServerError, err.Error()
	}

	return fiber.StatusCreated, credential
}
