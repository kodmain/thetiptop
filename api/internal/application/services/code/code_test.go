package code_test

import (
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/kodmain/thetiptop/api/internal/application/services/code"
	"github.com/kodmain/thetiptop/api/internal/domain/code/entities"
	"github.com/kodmain/thetiptop/api/internal/infrastructure/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// CodeServiceMock est le mock pour CodeServiceInterface
type CodeServiceMock struct {
	mock.Mock
}

func (m *CodeServiceMock) ListErrors() (map[string]*entities.Code, errors.ErrorInterface) {
	args := m.Called()
	if args.Get(0) != nil && args.Get(1) == nil {
		return args.Get(0).(map[string]*entities.Code), nil
	}
	return nil, args.Get(1).(errors.ErrorInterface)
}

func TestListErrors_Success(t *testing.T) {
	// Initialiser le mock
	mockService := new(CodeServiceMock)

	// Préparer les données mock
	mockErrorMap := map[string]*entities.Code{
		"Error1": {
			Code:    400,
			Message: "Bad Request",
		},
		"Error2": {
			Code:    404,
			Message: "Not Found",
		},
	}

	// Définir les attentes sur le mock
	mockService.On("ListErrors").Return(mockErrorMap, nil)

	// Appel de la fonction à tester
	status, result := code.ListErrors(mockService)

	// Vérification du statut et des résultats
	assert.Equal(t, fiber.StatusOK, status)
	assert.NotNil(t, result)
	assert.Equal(t, mockErrorMap, result)

	// Vérifier que toutes les attentes du mock sont satisfaites
	mockService.AssertExpectations(t)
}

func TestListErrors_Error(t *testing.T) {
	// Initialiser le mock
	mockService := new(CodeServiceMock)

	// Préparer une erreur mock
	mockError := errors.ErrInternalServer

	// Définir les attentes sur le mock
	mockService.On("ListErrors").Return(nil, mockError)

	// Appel de la fonction à tester
	status, result := code.ListErrors(mockService)

	// Vérification du statut et des résultats
	assert.Equal(t, mockError.Code(), status)
	assert.NotNil(t, result) // Vérifie que le résultat est bien nil lorsque l'erreur survient

	// Vérifier que toutes les attentes du mock sont satisfaites
	mockService.AssertExpectations(t)
}
