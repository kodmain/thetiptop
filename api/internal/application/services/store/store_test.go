package services_test

import (
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/gofiber/fiber/v2"
	services "github.com/kodmain/thetiptop/api/internal/application/services/store"
	transfert "github.com/kodmain/thetiptop/api/internal/application/transfert/crm"
	"github.com/kodmain/thetiptop/api/internal/domain/store/entities"
	errors_domain_store "github.com/kodmain/thetiptop/api/internal/domain/store/errors"
	"github.com/kodmain/thetiptop/api/internal/infrastructure/errors"
	"github.com/stretchr/testify/assert"
)

// TestListStores teste la fonction ListStores du package store
func TestListStores(t *testing.T) {
	t.Run("successful list", func(t *testing.T) {
		mockService, cleanup := setup()
		defer cleanup()

		// Configuration du mock pour retourner une liste de magasins
		expectedStores := []*entities.Store{
			{ID: "store-123", Label: aws.String("Store One"), IsOnline: aws.Bool(true)},
			{ID: "store-456", Label: aws.String("Store Two"), IsOnline: aws.Bool(false)},
		}
		mockService.On("ListStores").Return(expectedStores, nil)

		// Appel de la fonction
		statusCode, response := services.ListStores(mockService)

		// Assertions
		assert.Equal(t, fiber.StatusOK, statusCode)
		assert.Equal(t, expectedStores, response)
		mockService.AssertExpectations(t)
	})

	t.Run("service error - internal server error", func(t *testing.T) {
		mockService, cleanup := setup()
		defer cleanup()

		// Configuration du mock pour retourner une erreur interne
		mockService.On("ListStores").Return(nil, errors.ErrInternalServer)

		// Appel de la fonction
		statusCode, response := services.ListStores(mockService)

		// Assertions
		assert.Equal(t, 500, statusCode)
		assert.Equal(t, errors.ErrInternalServer, response)
		mockService.AssertExpectations(t)
	})
}

func TestGetStoreByID(t *testing.T) {
	t.Run("validation error - missing ID", func(t *testing.T) {
		mockService, cleanup := setup()
		defer cleanup()

		// DTO avec ID manquant
		dto := &transfert.Store{ID: nil}

		// Appel de la fonction
		statusCode, response := services.GetStoreByID(mockService, dto)

		// Assertions
		assert.Equal(t, fiber.StatusBadRequest, statusCode)
		assert.Error(t, response.(errors.Errors))
	})

	t.Run("validation error - invalid ID format", func(t *testing.T) {
		mockService, cleanup := setup()
		defer cleanup()

		// DTO avec ID invalide
		invalidID := "invalid-uuid"
		dto := &transfert.Store{ID: &invalidID}

		// Appel de la fonction
		statusCode, response := services.GetStoreByID(mockService, dto)

		// Assertions
		assert.Equal(t, fiber.StatusBadRequest, statusCode)
		assert.Error(t, response.(errors.Errors))
	})

	t.Run("service error - store not found", func(t *testing.T) {
		mockService, cleanup := setup()
		defer cleanup()

		// DTO valide avec UUID
		storeID := "42debee6-2063-4566-baf1-37a7bdd139ff" // UUID valide
		dto := &transfert.Store{ID: &storeID}

		// Configuration du mock pour renvoyer une erreur de magasin non trouvé
		mockService.On("GetStoreByID", dto).Return(nil, errors_domain_store.ErrStoreNotFound)

		// Appel de la fonction
		statusCode, response := services.GetStoreByID(mockService, dto)

		// Assertions
		assert.Equal(t, fiber.StatusNotFound, statusCode)
		assert.Equal(t, errors_domain_store.ErrStoreNotFound, response)
		mockService.AssertExpectations(t)
	})

	t.Run("service error - internal server error", func(t *testing.T) {
		mockService, cleanup := setup()
		defer cleanup()

		// DTO valide avec UUID
		storeID := "42debee6-2063-4566-baf1-37a7bdd139ff" // UUID valide
		dto := &transfert.Store{ID: &storeID}

		// Configuration du mock pour renvoyer une erreur interne
		mockService.On("GetStoreByID", dto).Return(nil, errors.ErrInternalServer)

		// Appel de la fonction
		statusCode, response := services.GetStoreByID(mockService, dto)

		// Assertions
		assert.Equal(t, fiber.StatusInternalServerError, statusCode)
		assert.Equal(t, errors.ErrInternalServer, response)
		mockService.AssertExpectations(t)
	})

	t.Run("successful get", func(t *testing.T) {
		mockService, cleanup := setup()
		defer cleanup()

		// DTO valide avec UUID
		storeID := "42debee6-2063-4566-baf1-37a7bdd139ff" // UUID valide
		dto := &transfert.Store{ID: &storeID}

		// Configuration du mock pour renvoyer un magasin
		expectedStore := &entities.Store{ID: storeID, Label: aws.String("Store One"), IsOnline: aws.Bool(true)}
		mockService.On("GetStoreByID", dto).Return(expectedStore, nil)

		// Appel de la fonction
		statusCode, response := services.GetStoreByID(mockService, dto)

		// Assertions
		assert.Equal(t, fiber.StatusOK, statusCode)
		assert.Equal(t, expectedStore, response)
		mockService.AssertExpectations(t)
	})
}
