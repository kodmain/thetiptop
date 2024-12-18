package services_test

import (
	"testing"

	"github.com/gofiber/fiber/v2"
	services "github.com/kodmain/thetiptop/api/internal/application/services/store"
	transfert "github.com/kodmain/thetiptop/api/internal/application/transfert/crm"
	"github.com/kodmain/thetiptop/api/internal/domain/store/entities"
	errors_domain_store "github.com/kodmain/thetiptop/api/internal/domain/store/errors"
	"github.com/kodmain/thetiptop/api/internal/infrastructure/errors"
	"github.com/stretchr/testify/assert"
)

// TestGetCaisse teste la fonction GetCaisse du package store
//
// Parameters:
// - t: *testing.T
//
// Returns:
// - None: aucun retour de valeur
// TestGetCaisse teste la fonction GetCaisse du package store
//
// Parameters:
// - t: *testing.T
//
// Returns:
// - None: aucun retour de valeur
func TestGetCaisse(t *testing.T) {
	t.Run("validation error - missing ID", func(t *testing.T) {
		mockService, cleanup := setup()
		defer cleanup()

		// DTO avec ID manquant
		dto := &transfert.Caisse{ID: nil}

		// Appel de la fonction
		statusCode, response := services.GetCaisse(mockService, dto)

		// Assertions
		assert.Equal(t, fiber.StatusBadRequest, statusCode)
		assert.Error(t, response.(*errors.Error))
	})

	t.Run("service error - caisse not found", func(t *testing.T) {
		mockService, cleanup := setup()
		defer cleanup()

		// DTO valide avec UUID
		caisseID := "42debee6-2063-4566-baf1-37a7bdd139ff" // UUID valide
		dto := &transfert.Caisse{ID: &caisseID}

		// Configuration du mock pour renvoyer une erreur de caisse non trouvée
		mockService.On("GetCaisse", dto).Return(nil, errors_domain_store.ErrCaisseNotFound)

		// Appel de la fonction
		statusCode, response := services.GetCaisse(mockService, dto)

		// Assertions
		assert.Equal(t, fiber.StatusNotFound, statusCode)
		assert.Equal(t, errors_domain_store.ErrCaisseNotFound, response)
		mockService.AssertExpectations(t)
	})

	t.Run("service error - internal server error", func(t *testing.T) {
		mockService, cleanup := setup()
		defer cleanup()

		// DTO valide avec UUID
		caisseID := "42debee6-2063-4566-baf1-37a7bdd139ff" // UUID valide
		dto := &transfert.Caisse{ID: &caisseID}

		// Configuration du mock pour renvoyer une erreur interne
		mockService.On("GetCaisse", dto).Return(nil, errors.ErrInternalServer)

		// Appel de la fonction
		statusCode, response := services.GetCaisse(mockService, dto)

		// Assertions
		assert.Equal(t, fiber.StatusInternalServerError, statusCode)
		assert.Equal(t, errors.ErrInternalServer, response)
		mockService.AssertExpectations(t)
	})

	t.Run("successful get", func(t *testing.T) {
		mockService, cleanup := setup()
		defer cleanup()

		// DTO valide avec UUID
		caisseID := "42debee6-2063-4566-baf1-37a7bdd139ff" // UUID valide
		storeID := "09054592-99a8-41fa-ab8c-6d9a70e49b64"
		dto := &transfert.Caisse{ID: &caisseID}

		// Configuration du mock pour renvoyer une caisse
		expectedCaisse := &entities.Caisse{ID: caisseID, StoreID: &storeID}
		mockService.On("GetCaisse", dto).Return(expectedCaisse, nil)

		// Appel de la fonction
		statusCode, response := services.GetCaisse(mockService, dto)

		// Assertions
		assert.Equal(t, fiber.StatusOK, statusCode)
		assert.Equal(t, expectedCaisse, response)
		mockService.AssertExpectations(t)
	})
}

// TestCreateCaisse teste la fonction CreateCaisse du package store
//
// Parameters:
// - t: *testing.T
//
// Returns:
// - None: aucun retour de valeur
// TestCreateCaisse teste la fonction CreateCaisse du package store
//
// Parameters:
// - t: *testing.T
//
// Returns:
// - None: aucun retour de valeur
func TestCreateCaisse(t *testing.T) {

	t.Run("validation error - missing store_id", func(t *testing.T) {
		mockService, cleanup := setup()
		defer cleanup()
		// DTO avec store_id manquant
		dto := &transfert.Caisse{StoreID: nil}

		// Appel de la fonction
		statusCode, response := services.CreateCaisse(mockService, dto)

		// Assertions
		assert.Equal(t, fiber.StatusBadRequest, statusCode)
		assert.Error(t, response.(*errors.Error))
	})

	t.Run("service error - creation failed", func(t *testing.T) {
		mockService, cleanup := setup()
		defer cleanup()
		// DTO valide
		storeID := "09054592-99a8-41fa-ab8c-6d9a70e49b64"
		dto := &transfert.Caisse{StoreID: &storeID}

		// Configuration du mock pour renvoyer une erreur de création
		mockService.On("CreateCaisse", dto).Return(nil, errors.ErrInternalServer)

		// Appel de la fonction
		statusCode, response := services.CreateCaisse(mockService, dto)

		// Assertions
		assert.Equal(t, fiber.StatusInternalServerError, statusCode)
		assert.Equal(t, errors.ErrInternalServer, response)
		mockService.AssertExpectations(t)
	})

	t.Run("successful creation", func(t *testing.T) {
		mockService, cleanup := setup()
		defer cleanup()
		// DTO valide
		caisseID := "09054592-99a8-41fa-ab8c-6d9a70e49b64"
		storeID := "09054592-99a8-41fa-ab8c-6d9a70e49b64"
		dto := &transfert.Caisse{StoreID: &storeID}

		// Configuration du mock pour renvoyer une caisse créée
		createdCaisse := &entities.Caisse{ID: caisseID, StoreID: &storeID}
		mockService.On("CreateCaisse", dto).Return(createdCaisse, nil)

		// Appel de la fonction
		statusCode, response := services.CreateCaisse(mockService, dto)

		// Assertions
		assert.Equal(t, fiber.StatusCreated, statusCode)
		assert.Equal(t, createdCaisse, response)
		mockService.AssertExpectations(t)
	})

	t.Run("successful get", func(t *testing.T) {
		mockService, cleanup := setup()
		defer cleanup()

		// DTO valide avec UUID
		caisseID := "42debee6-2063-4566-baf1-37a7bdd139ff" // UUID valide
		storeID := "09054592-99a8-41fa-ab8c-6d9a70e49b64"
		dto := &transfert.Caisse{ID: &caisseID}

		// Configuration du mock pour renvoyer une caisse
		expectedCaisse := &entities.Caisse{ID: caisseID, StoreID: &storeID}
		mockService.On("GetCaisse", dto).Return(expectedCaisse, nil)

		// Appel de la fonction
		statusCode, response := services.GetCaisse(mockService, dto)

		// Assertions
		assert.Equal(t, fiber.StatusOK, statusCode)
		assert.Equal(t, expectedCaisse, response)
		mockService.AssertExpectations(t)
	})
}

// TestDeleteCaisse teste la fonction DeleteCaisse du package store
//
// Parameters:
// - t: *testing.T
//
// Returns:
// - None: aucun retour de valeur
func TestDeleteCaisse(t *testing.T) {
	t.Run("validation error - missing ID", func(t *testing.T) {
		mockService, cleanup := setup()
		defer cleanup()
		// DTO avec ID manquant
		dto := &transfert.Caisse{ID: nil}

		// Appel de la fonction
		statusCode, response := services.DeleteCaisse(mockService, dto)

		// Assertions
		assert.Equal(t, fiber.StatusBadRequest, statusCode)
		assert.Error(t, response.(*errors.Error))
	})

	t.Run("service error - caisse not found", func(t *testing.T) {
		mockService, cleanup := setup()
		defer cleanup()
		// DTO valide
		caisseID := "09054592-99a8-41fa-ab8c-6d9a70e49b64"
		dto := &transfert.Caisse{ID: &caisseID}

		// Configuration du mock pour renvoyer une erreur de caisse non trouvée
		mockService.On("DeleteCaisse", dto).Return(errors_domain_store.ErrCaisseNotFound)

		// Appel de la fonction
		statusCode, response := services.DeleteCaisse(mockService, dto)

		// Assertions
		assert.Equal(t, fiber.StatusNotFound, statusCode)
		assert.Equal(t, errors_domain_store.ErrCaisseNotFound, response)
		mockService.AssertExpectations(t)
	})

	t.Run("service error - internal server error", func(t *testing.T) {
		mockService, cleanup := setup()
		defer cleanup()
		// DTO valide
		caisseID := "09054592-99a8-41fa-ab8c-6d9a70e49b64"
		dto := &transfert.Caisse{ID: &caisseID}

		// Configuration du mock pour renvoyer une erreur interne
		mockService.On("DeleteCaisse", dto).Return(errors.ErrInternalServer)

		// Appel de la fonction
		statusCode, response := services.DeleteCaisse(mockService, dto)

		// Assertions
		assert.Equal(t, fiber.StatusInternalServerError, statusCode)
		assert.Equal(t, errors.ErrInternalServer, response)
		mockService.AssertExpectations(t)
	})

	t.Run("successful deletion", func(t *testing.T) {
		mockService, cleanup := setup()
		defer cleanup()
		// DTO valide
		caisseID := "09054592-99a8-41fa-ab8c-6d9a70e49b64"
		dto := &transfert.Caisse{ID: &caisseID}

		// Configuration du mock pour réussir la suppression
		mockService.On("DeleteCaisse", dto).Return(nil)

		// Appel de la fonction
		statusCode, response := services.DeleteCaisse(mockService, dto)

		// Assertions
		assert.Equal(t, fiber.StatusNoContent, statusCode)
		assert.Nil(t, response)
		mockService.AssertExpectations(t)
	})
}

// TestUpdateCaisse teste la fonction UpdateCaisse du package store
//
// Parameters:
// - t: *testing.T
//
// Returns:
// - None: aucun retour de valeur
func TestUpdateCaisse(t *testing.T) {
	t.Run("validation error - missing ID", func(t *testing.T) {
		mockService, cleanup := setup()
		defer cleanup()
		// DTO avec ID manquant
		storeID := "09054592-99a8-41fa-ab8c-6d9a70e49b64"
		dto := &transfert.Caisse{ID: nil, StoreID: &storeID}

		// Appel de la fonction
		statusCode, response := services.UpdateCaisse(mockService, dto)

		// Assertions
		assert.Equal(t, fiber.StatusBadRequest, statusCode)
		assert.Error(t, response.(*errors.Error))
	})

	t.Run("validation error - missing store_id", func(t *testing.T) {
		mockService, cleanup := setup()
		defer cleanup()
		// DTO avec store_id manquant
		caisseID := "09054592-99a8-41fa-ab8c-6d9a70e49b64"
		dto := &transfert.Caisse{ID: &caisseID, StoreID: nil}

		// Appel de la fonction
		statusCode, response := services.UpdateCaisse(mockService, dto)

		// Assertions
		assert.Equal(t, fiber.StatusBadRequest, statusCode)
		assert.Error(t, response.(*errors.Error))
	})

	t.Run("service error - caisse not found", func(t *testing.T) {
		mockService, cleanup := setup()
		defer cleanup()
		// DTO valide
		caisseID := "09054592-99a8-41fa-ab8c-6d9a70e49b64"
		storeID := "09054592-99a8-41fa-ab8c-6d9a70e49b64"
		dto := &transfert.Caisse{ID: &caisseID, StoreID: &storeID}

		// Configuration du mock pour renvoyer une erreur de caisse non trouvée
		mockService.On("UpdateCaisse", dto).Return(nil, errors_domain_store.ErrCaisseNotFound)

		// Appel de la fonction
		statusCode, response := services.UpdateCaisse(mockService, dto)

		// Assertions
		assert.Equal(t, fiber.StatusNotFound, statusCode)
		assert.Equal(t, errors_domain_store.ErrCaisseNotFound, response)
		mockService.AssertExpectations(t)
	})

	t.Run("service error - internal server error", func(t *testing.T) {
		mockService, cleanup := setup()
		defer cleanup()
		// DTO valide
		caisseID := "09054592-99a8-41fa-ab8c-6d9a70e49b64"
		storeID := "09054592-99a8-41fa-ab8c-6d9a70e49b64"
		dto := &transfert.Caisse{ID: &caisseID, StoreID: &storeID}

		// Configuration du mock pour renvoyer une erreur interne
		mockService.On("UpdateCaisse", dto).Return(nil, errors.ErrInternalServer)

		// Appel de la fonction
		statusCode, response := services.UpdateCaisse(mockService, dto)

		// Assertions
		assert.Equal(t, fiber.StatusInternalServerError, statusCode)
		assert.Equal(t, errors.ErrInternalServer, response)
		mockService.AssertExpectations(t)
	})

	t.Run("successful update", func(t *testing.T) {
		mockService, cleanup := setup()
		defer cleanup()
		// DTO valide
		caisseID := "09054592-99a8-41fa-ab8c-6d9a70e49b64"
		storeID := "09054592-99a8-41fa-ab8c-6d9a70e49b64"
		dto := &transfert.Caisse{ID: &caisseID, StoreID: &storeID}

		// Configuration du mock pour renvoyer une caisse mise à jour
		updatedCaisse := &entities.Caisse{ID: caisseID, StoreID: &storeID}
		mockService.On("UpdateCaisse", dto).Return(updatedCaisse, nil)

		// Appel de la fonction
		statusCode, response := services.UpdateCaisse(mockService, dto)

		// Assertions
		assert.Equal(t, fiber.StatusOK, statusCode)
		assert.Equal(t, updatedCaisse, response)
		mockService.AssertExpectations(t)
	})
}
