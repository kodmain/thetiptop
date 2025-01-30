package services_test

import (
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/gofiber/fiber/v2"
	services "github.com/kodmain/thetiptop/api/internal/application/services/user"
	transfert "github.com/kodmain/thetiptop/api/internal/application/transfert/user"
	"github.com/kodmain/thetiptop/api/internal/domain/user/entities"
	errors_domain_user "github.com/kodmain/thetiptop/api/internal/domain/user/errors"
	"github.com/kodmain/thetiptop/api/internal/infrastructure/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestRegisterEmployee(t *testing.T) {
	// Variables (si besoin) pour email / password
	email := "test@example.com"
	password := "ValidP@ssw0rd"
	passwordSyntaxFail := "short"

	t.Run("invalid password", func(t *testing.T) {
		t.Parallel()

		mockService := new(DomainUserService)
		// Pas d'appel effectif à RegisterEmployee attendu
		// car la validation échoue avant. On peut soit ne pas mettre .On(...),
		// soit mettre On(...).Maybe() pour que mock n'exige pas forcément l'appel.

		statusCode, response := services.RegisterEmployee(
			mockService,
			&transfert.Credential{
				Email:    &email,
				Password: &passwordSyntaxFail,
			},
			&transfert.Employee{},
		)

		assert.Equal(t, fiber.StatusBadRequest, statusCode)
		assert.NotNil(t, response)

		// Vérifie qu'on n'a PAS appelé RegisterEmployee
		mockService.AssertNotCalled(t, "RegisterEmployee", mock.AnythingOfType("*transfert.Credential"), mock.AnythingOfType("*transfert.Employee"))
	})

	t.Run("valid password and fields", func(t *testing.T) {
		t.Parallel()

		mockService := new(DomainUserService)
		mockService.On("RegisterEmployee", mock.AnythingOfType("*transfert.Credential"), mock.AnythingOfType("*transfert.Employee")).
			Return(&entities.Employee{}, nil)

		statusCode, response := services.RegisterEmployee(
			mockService,
			&transfert.Credential{
				Email:    &email,
				Password: &password,
			},
			&transfert.Employee{},
		)

		assert.Equal(t, fiber.StatusCreated, statusCode)
		assert.NotNil(t, response)

		mockService.AssertExpectations(t)
	})

	t.Run("employee already exists", func(t *testing.T) {
		t.Parallel()

		mockService := new(DomainUserService)
		mockService.On("RegisterEmployee", mock.AnythingOfType("*transfert.Credential"), mock.AnythingOfType("*transfert.Employee")).
			Return(nil, errors_domain_user.ErrCredentialAlreadyExists)

		statusCode, response := services.RegisterEmployee(
			mockService,
			&transfert.Credential{
				Email:    &email,
				Password: &password,
			},
			&transfert.Employee{},
		)

		assert.Equal(t, fiber.StatusConflict, statusCode)
		assert.NotNil(t, response)

		mockService.AssertExpectations(t)
	})

	t.Run("server error during registration", func(t *testing.T) {
		t.Parallel()

		mockService := new(DomainUserService)
		mockService.On("RegisterEmployee", mock.AnythingOfType("*transfert.Credential"), mock.AnythingOfType("*transfert.Employee")).
			Return(nil, errors.ErrInternalServer)

		statusCode, response := services.RegisterEmployee(
			mockService,
			&transfert.Credential{
				Email:    &email,
				Password: &password,
			},
			&transfert.Employee{},
		)

		assert.Equal(t, fiber.StatusInternalServerError, statusCode)
		errObj, ok := response.(*errors.Error)
		assert.True(t, ok, "response should be of type *errors.Error")
		if ok {
			assert.Error(t, errObj)
		}

		mockService.AssertExpectations(t)
	})
}

func TestDeleteEmployee(t *testing.T) {
	t.Run("should return 400 if validation fails", func(t *testing.T) {
		t.Parallel()

		mockService := new(DomainUserService)
		dtoEmployee := &transfert.Employee{ID: nil}

		statusCode, response := services.DeleteEmployee(mockService, dtoEmployee)

		assert.Equal(t, fiber.StatusBadRequest, statusCode)

		errMaps, ok := response.(errors.Errors)
		assert.True(t, ok, "response should be of type errors.Errors")
		if ok {
			assert.Error(t, errMaps)
		}

		// On ne s'attend pas à un appel, la validation échoue
		mockService.AssertNotCalled(t, "DeleteEmployee", mock.Anything)
	})

	t.Run("should return 404 if employee not found", func(t *testing.T) {
		t.Parallel()

		mockService := new(DomainUserService)
		employeeID := "123e4567-e89b-12d3-a456-426614174000"
		dtoEmployee := &transfert.Employee{ID: &employeeID}

		mockService.On("DeleteEmployee", dtoEmployee).Return(errors_domain_user.ErrEmployeeNotFound)

		statusCode, response := services.DeleteEmployee(mockService, dtoEmployee)

		assert.Equal(t, fiber.StatusNotFound, statusCode)
		assert.Equal(t, errors_domain_user.ErrEmployeeNotFound, response)

		mockService.AssertExpectations(t)
	})

	t.Run("should return 500 if internal server error occurs", func(t *testing.T) {
		t.Parallel()

		mockService := new(DomainUserService)
		employeeID := "123e4567-e89b-12d3-a456-426614174000"
		dtoEmployee := &transfert.Employee{ID: &employeeID}

		mockService.On("DeleteEmployee", dtoEmployee).Return(errors.ErrInternalServer)

		statusCode, response := services.DeleteEmployee(mockService, dtoEmployee)

		assert.Equal(t, fiber.StatusInternalServerError, statusCode)

		errObj, ok := response.(*errors.Error)
		assert.True(t, ok, "response should be of type *errors.Error")
		if ok {
			assert.Error(t, errObj)
		}

		mockService.AssertExpectations(t)
	})

	t.Run("should return 204 if employee is deleted successfully", func(t *testing.T) {
		t.Parallel()

		mockService := new(DomainUserService)
		employeeID := "123e4567-e89b-12d3-a456-426614174000"
		dtoEmployee := &transfert.Employee{ID: &employeeID}

		mockService.On("DeleteEmployee", dtoEmployee).Return(nil)

		statusCode, response := services.DeleteEmployee(mockService, dtoEmployee)

		assert.Equal(t, fiber.StatusNoContent, statusCode)
		assert.Nil(t, response)

		mockService.AssertExpectations(t)
	})
}

func TestUpdateEmployee(t *testing.T) {
	t.Run("invalid employee data", func(t *testing.T) {
		t.Parallel()

		mockService := new(DomainUserService)
		// ID manquant
		statusCode, response := services.UpdateEmployee(mockService, &transfert.Employee{
			ID: nil,
		})

		assert.Equal(t, fiber.StatusBadRequest, statusCode)
		errMaps, ok := response.(errors.Errors)
		assert.True(t, ok, "response should be of type errors.Errors")
		if ok {
			assert.Error(t, errMaps)
		}

		mockService.AssertNotCalled(t, "UpdateEmployee", mock.Anything)
	})

	t.Run("successful employee update", func(t *testing.T) {
		t.Parallel()

		mockService := new(DomainUserService)
		mockService.On("UpdateEmployee", mock.AnythingOfType("*transfert.Employee")).
			Return(&entities.Employee{
				ID: "123e4567-e89b-12d3-a456-426614174000",
			}, nil)

		statusCode, response := services.UpdateEmployee(mockService, &transfert.Employee{
			ID: aws.String("123e4567-e89b-12d3-a456-426614174000"),
		})

		assert.Equal(t, fiber.StatusOK, statusCode)
		assert.NotNil(t, response)

		mockService.AssertExpectations(t)
	})

	t.Run("employee update error", func(t *testing.T) {
		t.Parallel()

		mockService := new(DomainUserService)
		mockService.On("UpdateEmployee", mock.AnythingOfType("*transfert.Employee")).
			Return(nil, errors.ErrInternalServer)

		statusCode, response := services.UpdateEmployee(mockService, &transfert.Employee{
			ID: aws.String("123e4567-e89b-12d3-a456-426614174000"),
		})

		assert.Equal(t, fiber.StatusInternalServerError, statusCode)
		errObj, ok := response.(*errors.Error)
		assert.True(t, ok, "response should be of type *errors.Error")
		if ok {
			assert.Error(t, errObj)
		}

		mockService.AssertExpectations(t)
	})
}

func TestGetEmployee(t *testing.T) {
	t.Run("validation error", func(t *testing.T) {
		t.Parallel()

		mockService := new(DomainUserService)
		dtoEmployee := &transfert.Employee{ID: nil}

		statusCode, response := services.GetEmployee(mockService, dtoEmployee)

		assert.Equal(t, fiber.StatusBadRequest, statusCode)
		errMaps, ok := response.(errors.Errors)
		assert.True(t, ok, "response should be of type errors.Errors")
		if ok {
			assert.Contains(t, errMaps, "id", "id field should be present")
		}

		mockService.AssertNotCalled(t, "GetEmployee", dtoEmployee)
	})

	t.Run("employee not found", func(t *testing.T) {
		t.Parallel()

		mockService := new(DomainUserService)
		dtoEmployee := &transfert.Employee{
			ID: aws.String("42debee6-2063-4566-baf1-37a7bdd139ff"),
		}

		mockService.On("GetEmployee", dtoEmployee).Return(nil, errors_domain_user.ErrEmployeeNotFound)

		statusCode, response := services.GetEmployee(mockService, dtoEmployee)

		assert.Equal(t, fiber.StatusNotFound, statusCode)
		assert.Equal(t, errors_domain_user.ErrEmployeeNotFound, response)

		mockService.AssertExpectations(t)
	})

	t.Run("employee random error", func(t *testing.T) {
		t.Parallel()

		mockService := new(DomainUserService)
		dtoEmployee := &transfert.Employee{
			ID: aws.String("42debee6-2063-4566-baf1-37a7bdd139ff"),
		}

		mockService.On("GetEmployee", dtoEmployee).Return(nil, errors.ErrInternalServer)

		statusCode, response := services.GetEmployee(mockService, dtoEmployee)

		assert.Equal(t, fiber.StatusInternalServerError, statusCode)

		errObj, ok := response.(*errors.Error)
		assert.True(t, ok, "response should be of type *errors.Error")
		if ok {
			assert.Error(t, errObj)
		}

		mockService.AssertExpectations(t)
	})

	t.Run("successful employee retrieval", func(t *testing.T) {
		t.Parallel()

		mockService := new(DomainUserService)
		dtoEmployee := &transfert.Employee{
			ID: aws.String("42debee6-2063-4566-baf1-37a7bdd139ff"),
		}

		expectedEmployee := &entities.Employee{
			ID: "42debee6-2063-4566-baf1-37a7bdd139ff",
		}

		mockService.On("GetEmployee", dtoEmployee).Return(expectedEmployee, nil)

		statusCode, response := services.GetEmployee(mockService, dtoEmployee)

		assert.Equal(t, fiber.StatusOK, statusCode)
		assert.Equal(t, expectedEmployee, response)

		mockService.AssertExpectations(t)
	})
}
