package services_test

import (
	"fmt"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/google/uuid"
	"github.com/kodmain/thetiptop/api/internal/application/transfert"
	"github.com/kodmain/thetiptop/api/internal/domain/user/entities"
	"github.com/kodmain/thetiptop/api/internal/domain/user/errors"
	"github.com/kodmain/thetiptop/api/internal/infrastructure/security/token"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestEmployeeRegister(t *testing.T) {
	// Variables communes
	idEmployee, err := uuid.Parse("42debee6-2063-4566-baf1-37a7bdd139ff")
	assert.NoError(t, err)
	idValidation, err := uuid.Parse("42debee6-2061-4566-baf1-37a7bdd139ff")
	assert.NoError(t, err)

	sidEmployee := idEmployee.String()

	inputEmployee := &transfert.Employee{}

	expectedEmployee := &entities.Employee{
		ID: idEmployee.String(),
		Validations: []*entities.Validation{
			{
				ID:         idValidation.String(),
				Token:      token.NewLuhn("666666").Pointer(),
				Type:       0,
				Validated:  false,
				EmployeeID: &sidEmployee,
			},
		},
	}

	inputCredential := &transfert.Credential{
		Email:    aws.String("hello@thetiptop"),
		Password: aws.String("azertyuiop"),
	}

	expectedCredential := &entities.Credential{
		ID:       idEmployee.String(),
		Email:    inputCredential.Email,
		Password: aws.String("$2a$10$wO5PfDAGp6w2ubKp0vEdXeUe2HlfOv5iRJ3C3MVR0vJhscD0G.NKS"), // hashed password
		//EmployeeID: &sidEmployee,
	}

	t.Run("nil input", func(t *testing.T) {
		service, _, _ := setup()
		require.NotNil(t, service)

		result, err := service.RegisterEmployee(nil, nil)
		require.Error(t, err)
		require.Nil(t, result)
		require.Equal(t, errors.ErrNoDto, err.Error())
	})

	t.Run("employee already exists", func(t *testing.T) {
		service, mockRepo, _ := setup()
		dtoCredential := &transfert.Credential{Email: aws.String("existing@example.com")}
		dtoEmployee := &transfert.Employee{}

		mockRepo.On("ReadCredential", dtoCredential).Return(&entities.Credential{}, nil)

		Employee, err := service.RegisterEmployee(dtoCredential, dtoEmployee)
		assert.Nil(t, Employee)
		assert.EqualError(t, err, errors.ErrEmployeeAlreadyExists)
		mockRepo.AssertExpectations(t)
	})

	t.Run("credential creation error", func(t *testing.T) {
		service, mockRepo, _ := setup()
		dtoCredential := &transfert.Credential{Email: aws.String("new@example.com")}
		dtoEmployee := &transfert.Employee{}

		mockRepo.On("ReadCredential", dtoCredential).Return(nil, fmt.Errorf("not found"))
		mockRepo.On("CreateCredential", dtoCredential).Return(nil, fmt.Errorf("error creating credential"))

		Employee, err := service.RegisterEmployee(dtoCredential, dtoEmployee)
		assert.Nil(t, Employee)
		assert.EqualError(t, err, "error creating credential")
		mockRepo.AssertExpectations(t)
	})

	t.Run("employee creation error", func(t *testing.T) {
		service, mockRepo, _ := setup()
		dtoCredential := &transfert.Credential{Email: aws.String("new@example.com")}
		dtoEmployee := &transfert.Employee{}

		mockRepo.On("ReadCredential", dtoCredential).Return(nil, fmt.Errorf("not found"))
		mockRepo.On("CreateCredential", dtoCredential).Return(expectedCredential, nil)
		mockRepo.On("CreateEmployee", dtoEmployee).Return(nil, fmt.Errorf("error creating employee"))

		Employee, err := service.RegisterEmployee(dtoCredential, dtoEmployee)
		assert.Nil(t, Employee)
		assert.EqualError(t, err, "error creating Employee")
		mockRepo.AssertExpectations(t)
	})

	t.Run("employee update error", func(t *testing.T) {
		service, mockRepo, _ := setup()

		mockRepo.On("ReadCredential", inputCredential).Return(nil, fmt.Errorf("not found"))
		mockRepo.On("CreateCredential", inputCredential).Return(expectedCredential, nil)
		mockRepo.On("CreateEmployee", inputEmployee).Return(expectedEmployee, nil)
		mockRepo.On("UpdateEmployee", expectedEmployee).Return(fmt.Errorf("error updating Employee"))

		Employee, err := service.RegisterEmployee(inputCredential, inputEmployee)
		assert.Nil(t, Employee)
		assert.EqualError(t, err, "error updating Employee")
		mockRepo.AssertExpectations(t)
	})

	t.Run("credential update error", func(t *testing.T) {
		service, mockRepo, _ := setup()

		mockRepo.On("ReadCredential", inputCredential).Return(nil, fmt.Errorf("not found"))
		mockRepo.On("CreateCredential", inputCredential).Return(expectedCredential, nil)
		mockRepo.On("CreateEmployee", inputEmployee).Return(expectedEmployee, nil)
		mockRepo.On("UpdateEmployee", expectedEmployee).Return(nil)
		mockRepo.On("UpdateCredential", expectedCredential).Return(fmt.Errorf("error updating credential"))

		Employee, err := service.RegisterEmployee(inputCredential, inputEmployee)
		assert.Nil(t, Employee)
		assert.EqualError(t, err, "error updating credential")
		mockRepo.AssertExpectations(t)
	})

	t.Run("successful employee and credential creation", func(t *testing.T) {
		service, mockRepo, mockMailer := setup()

		mockRepo.On("ReadCredential", inputCredential).Return(nil, fmt.Errorf("not found"))
		mockRepo.On("CreateCredential", inputCredential).Return(expectedCredential, nil)
		mockRepo.On("CreateEmployee", inputEmployee).Return(expectedEmployee, nil)
		mockRepo.On("UpdateEmployee", expectedEmployee).Return(nil)
		mockRepo.On("UpdateCredential", expectedCredential).Return(nil)

		mockMailer.On("Send", mock.AnythingOfType("*mail.Mail")).Return(nil)

		Employee, err := service.RegisterEmployee(inputCredential, inputEmployee)
		assert.NotNil(t, Employee)
		assert.NoError(t, err)
		assert.Equal(t, sidEmployee, Employee.ID)
		mockRepo.AssertExpectations(t)
	})
}

func TestGetEmployee(t *testing.T) {
	t.Run("employee not found", func(t *testing.T) {
		service, mockRepo, _ := setup()

		dtoEmployee := &transfert.Employee{ID: aws.String("employee-id")}

		mockRepo.On("ReadEmployee", dtoEmployee).Return(nil, fmt.Errorf(errors.ErrEmployeeNotFound))

		employee, err := service.GetEmployee(dtoEmployee)
		assert.Nil(t, employee)
		assert.EqualError(t, err, errors.ErrEmployeeNotFound)
		mockRepo.AssertExpectations(t)
	})

	t.Run("successful retrieval", func(t *testing.T) {
		service, mockRepo, _ := setup()

		dtoEmployee := &transfert.Employee{ID: aws.String("employee-id")}
		expectedEmployee := &entities.Employee{ID: "employee-id"}

		mockRepo.On("ReadEmployee", dtoEmployee).Return(expectedEmployee, nil)

		employee, err := service.GetEmployee(dtoEmployee)
		assert.NotNil(t, employee)
		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})
}

func TestDeleteEmployee(t *testing.T) {
	t.Run("nil employee DTO returns error", func(t *testing.T) {
		service, _, _ := setup()

		err := service.DeleteEmployee(nil)
		assert.EqualError(t, err, errors.ErrNoDto)
	})

	t.Run("successful deletion", func(t *testing.T) {
		service, mockRepo, _ := setup()

		dtoEmployee := &transfert.Employee{ID: aws.String("employee-id")}

		mockRepo.On("DeleteEmployee", dtoEmployee).Return(nil)

		err := service.DeleteEmployee(dtoEmployee)
		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("deletion error", func(t *testing.T) {
		service, mockRepo, _ := setup()

		dtoEmployee := &transfert.Employee{ID: aws.String("employee-id")}

		mockRepo.On("DeleteEmployee", dtoEmployee).Return(fmt.Errorf("deletion failed"))

		err := service.DeleteEmployee(dtoEmployee)
		assert.EqualError(t, err, "deletion failed")
		mockRepo.AssertExpectations(t)
	})
}

func TestUpdateEmployee(t *testing.T) {
	t.Run("employee not found", func(t *testing.T) {
		service, mockRepo, _ := setup()

		dtoEmployee := &transfert.Employee{ID: aws.String("employee-id")}

		mockRepo.On("ReadEmployee", mock.AnythingOfType("*transfert.Employee")).Return(nil, fmt.Errorf(errors.ErrEmployeeNotFound))

		employee, err := service.UpdateEmployee(dtoEmployee)
		assert.Nil(t, employee)
		assert.EqualError(t, err, errors.ErrEmployeeNotFound)
		mockRepo.AssertExpectations(t)
	})

	t.Run("successful update", func(t *testing.T) {
		service, mockRepo, _ := setup()

		dtoEmployee := &transfert.Employee{ID: aws.String("employee-id")}
		existingEmployee := &entities.Employee{ID: "employee-id"}

		mockRepo.On("ReadEmployee", mock.AnythingOfType("*transfert.Employee")).Return(existingEmployee, nil)
		mockRepo.On("UpdateEmployee", existingEmployee).Return(nil)

		employee, err := service.UpdateEmployee(dtoEmployee)
		assert.NotNil(t, employee)
		assert.NoError(t, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("update error", func(t *testing.T) {
		service, mockRepo, _ := setup()

		dtoEmployee := &transfert.Employee{ID: aws.String("employee-id")}
		existingEmployee := &entities.Employee{ID: "employee-id"}

		mockRepo.On("ReadEmployee", mock.AnythingOfType("*transfert.Employee")).Return(existingEmployee, nil)
		mockRepo.On("UpdateEmployee", existingEmployee).Return(fmt.Errorf("update failed"))

		employee, err := service.UpdateEmployee(dtoEmployee)
		assert.Nil(t, employee)
		assert.EqualError(t, err, "update failed")
		mockRepo.AssertExpectations(t)
	})
}
