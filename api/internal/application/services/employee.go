package services

import (
	"github.com/gofiber/fiber/v2"
	"github.com/kodmain/thetiptop/api/internal/application/transfert"
	"github.com/kodmain/thetiptop/api/internal/application/validator"
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
	if err := dtoEmployee.Check(data.Validator{
		"id": {validator.Required, validator.ID},
	}); err != nil {
		return err.Code(), err
	}

	// Attempt to delete the employee using the service
	if err := service.DeleteEmployee(dtoEmployee); err != nil {
		return err.Code(), err
	}

	// Return 204 if the deletion is successful with no content
	return fiber.StatusNoContent, nil
}

func GetEmployee(service services.UserServiceInterface, dtoEmployee *transfert.Employee) (int, any) {
	if err := dtoEmployee.Check(data.Validator{
		"id": {validator.Required, validator.ID},
	}); err != nil {
		return err.Code(), err
	}

	employee, err := service.GetEmployee(dtoEmployee)
	if err != nil {
		return err.Code(), err
	}

	return fiber.StatusOK, employee
}

func UpdateEmployee(service services.UserServiceInterface, employeeDTO *transfert.Employee) (int, any) {
	if err := employeeDTO.Check(data.Validator{
		"id": {validator.Required, validator.ID},
	}); err != nil {
		return err.Code(), err
	}

	employee, err := service.UpdateEmployee(employeeDTO)
	if err != nil {
		return err.Code(), err
	}

	return fiber.StatusOK, employee
}

func RegisterEmployee(service services.UserServiceInterface, credentialDTO *transfert.Credential, employeeDTO *transfert.Employee) (int, any) {
	if err := credentialDTO.Check(data.Validator{
		"email":    {validator.Required, validator.Email},
		"password": {validator.Required, validator.Password},
	}); err != nil {
		return err.Code(), err
	}

	if err := employeeDTO.Check(data.Validator{}); err != nil {
		return err.Code(), err
	}

	credential, err := service.RegisterEmployee(credentialDTO, employeeDTO)
	if err != nil {
		return err.Code(), err
	}

	return fiber.StatusCreated, credential
}
