package entities_test

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/kodmain/thetiptop/api/internal/application/transfert"
	"github.com/kodmain/thetiptop/api/internal/domain/user/entities"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func TestEmployee_HasSuccessValidation(t *testing.T) {
	employee := &entities.Employee{
		Validations: entities.Validations{
			&entities.Validation{
				Type:      entities.PasswordRecover,
				Validated: true,
			},
			&entities.Validation{
				Type:      entities.PasswordRecover,
				Validated: false,
			},
			&entities.Validation{
				Type:      entities.PasswordRecover,
				Validated: true,
			},
		},
	}

	validationType := entities.PasswordRecover
	result := employee.HasSuccessValidation(entities.PasswordRecover)

	assert.NotNil(t, result)
	assert.Equal(t, validationType, result.Type)
	assert.True(t, result.Validated)

	result = employee.HasSuccessValidation(entities.MailValidation)
	assert.Nil(t, result)
}

func TestEmployeeHasNotExpiredValidation(t *testing.T) {
	employee := &entities.Employee{
		Validations: entities.Validations{
			&entities.Validation{
				Type:      entities.PasswordRecover,
				Validated: true,
				ExpiresAt: time.Now().Add(time.Hour),
			},
			&entities.Validation{
				Type:      entities.PasswordRecover,
				Validated: false,
				ExpiresAt: time.Now().Add(-time.Hour),
			},
			&entities.Validation{
				Type:      entities.PasswordRecover,
				Validated: false,
				ExpiresAt: time.Now().Add(time.Hour),
			},
		},
	}

	validationType := entities.PasswordRecover
	result := employee.HasNotExpiredValidation(entities.PasswordRecover)

	assert.NotNil(t, result)
	assert.Equal(t, validationType, result.Type)
	assert.False(t, result.Validated)
	assert.False(t, result.HasExpired())

	result = employee.HasNotExpiredValidation(entities.MailValidation)
	assert.Nil(t, result)
}

func TestEmployeeBeforeCreateAndUpdate(t *testing.T) {
	employee := &entities.Employee{}

	employee.Validations = append(employee.Validations, &entities.Validation{
		Type: entities.MailValidation,
	})

	// Test BeforeCreate
	err := employee.BeforeCreate(nil)
	assert.Nil(t, err)
	assert.NotEmpty(t, employee.ID)
	assert.Equal(t, employee.ID, employee.ID)

	// Simule un timestamp précédent pour comparer avec AfterUpdate
	old := employee.UpdatedAt
	time.Sleep(100 * time.Millisecond)

	// Test BeforeUpdate
	err = employee.BeforeUpdate(nil)
	assert.Nil(t, err)
	assert.True(t, employee.UpdatedAt.After(old))
}

func TestEmployee_AfterFind(t *testing.T) {
	// Setup SQLite in-memory database for testing
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	assert.Nil(t, err)

	// Automigrate for Employee and Validation tables
	err = db.AutoMigrate(&entities.Employee{}, &entities.Validation{})
	assert.Nil(t, err)

	// Create a new employee with no validations
	employee := &entities.Employee{
		ID: uuid.New().String(),
	}

	// Insert the employee into the database
	err = db.Create(&employee).Error
	assert.Nil(t, err)

	// Call AfterFind by retrieving the employee from the database
	err = db.First(&employee, "id = ?", employee.ID).Error
	assert.Nil(t, err)

	// Test if the Validations relation is properly loaded
	assert.NotNil(t, employee.Validations)
	assert.True(t, len(employee.Validations) == 0) // The employee should have 0 validations
}

func TestCreateEmployee(t *testing.T) {
	// Crée un objet de type transfert.Employee avec les champs nécessaires
	input := &transfert.Employee{}

	// Appelle la fonction CreateEmployee
	employee := entities.CreateEmployee(input)

	// Vérifie que le employee n'est pas nil
	assert.NotNil(t, employee)

	// Vérifie que le champ Validations est bien initialisé à une liste vide
	assert.NotNil(t, employee.Validations)
	assert.Equal(t, 0, len(employee.Validations))
}

func TestCreateEmployeeWithNilFields(t *testing.T) {
	// Cas de test où CGU et Newsletter sont nil
	input := &transfert.Employee{}

	// Appelle la fonction CreateEmployee
	employee := entities.CreateEmployee(input)

	// Vérifie que le employee n'est pas nil
	assert.NotNil(t, employee)

	// Vérifie que le champ Validations est bien initialisé à une liste vide
	assert.NotNil(t, employee.Validations)
	assert.Equal(t, 0, len(employee.Validations))
}
