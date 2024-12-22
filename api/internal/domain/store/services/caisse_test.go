// services_test.go
//
// Description:
// Ce fichier contient les tests unitaires pour les méthodes GetCaisse, CreateCaisse, UpdateCaisse et DeleteCaisse du StoreService.
// Les scénarios de test sont en français, mais les commentaires de code sont en anglais.
// Les commentaires de code suivent le format demandé.
//
// Les modifications apportées incluent l'utilisation de mock.Anything pour gérer le deuxième paramètre options ...database.Option.
//
// Parameters:
// - None
//
// Returns:
// - None: no return value

package services_test

import (
	"testing"

	"github.com/kodmain/thetiptop/api/internal/application/security"
	transfert "github.com/kodmain/thetiptop/api/internal/application/transfert/crm"
	"github.com/kodmain/thetiptop/api/internal/domain/store/entities"
	user "github.com/kodmain/thetiptop/api/internal/domain/user/entities"
	"github.com/kodmain/thetiptop/api/internal/infrastructure/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Test_GetCaisse tests the GetCaisse method of StoreService
// Parameters:
// - t: *testing.T
//
// Returns:
// - None: no return value
func Test_GetCaisse(t *testing.T) {
	t.Run("Devrait retourner une caisse lorsque autorisé et existe", func(t *testing.T) {
		service, mockRepo, mockPerms := setup()

		idCaisse := "caisse-123"
		dto := &transfert.Caisse{ID: &idCaisse}

		caisse := &entities.Caisse{ID: "caisse-123"}

		mockPerms.On("IsGrantedByRoles", []security.Role{user.ROLE_EMPLOYEE}).Return(true)
		mockRepo.On("ReadCaisse", dto, mock.Anything).Return(caisse, nil)

		result, err := service.GetCaisse(dto)
		assert.Nil(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, caisse, result)

		mockRepo.AssertExpectations(t)
		mockPerms.AssertExpectations(t)
	})

	t.Run("Devrait retourner une erreur lorsque non autorisé", func(t *testing.T) {
		service, _, mockPerms := setup()

		idCaisse := "caisse-123"
		dto := &transfert.Caisse{ID: &idCaisse}

		mockPerms.On("IsGrantedByRoles", []security.Role{user.ROLE_EMPLOYEE}).Return(false)

		result, err := service.GetCaisse(dto)
		assert.Nil(t, result)
		assert.NotNil(t, err)
		assert.Equal(t, errors.ErrUnauthorized, err)

		mockPerms.AssertExpectations(t)
	})

	t.Run("Devrait retourner une erreur lorsque dto est nil", func(t *testing.T) {
		service, _, _ := setup()

		result, err := service.GetCaisse(nil)
		assert.Nil(t, result)
		assert.NotNil(t, err)
		assert.Equal(t, errors.ErrNoDto, err)
	})

	t.Run("Devrait retourner une erreur lorsque le repo échoue", func(t *testing.T) {
		service, mockRepo, mockPerms := setup()

		idCaisse := "caisse-123"
		dto := &transfert.Caisse{ID: &idCaisse}

		mockPerms.On("IsGrantedByRoles", []security.Role{user.ROLE_EMPLOYEE}).Return(true)
		mockRepo.On("ReadCaisse", dto, mock.Anything).Return(nil, errors.ErrNoData)

		result, err := service.GetCaisse(dto)
		assert.Nil(t, result)
		assert.NotNil(t, err)
		assert.Equal(t, errors.ErrNoData, err)

		mockRepo.AssertExpectations(t)
		mockPerms.AssertExpectations(t)
	})
}

// Test_CreateCaisse tests the CreateCaisse method of StoreService
// Parameters:
// - t: *testing.T
//
// Returns:
// - None: no return value
func Test_CreateCaisse(t *testing.T) {
	t.Run("Devrait créer une caisse lorsque autorisé et store existe", func(t *testing.T) {
		service, mockRepo, mockPerms := setup()

		idCaisse := "caisse-123"
		idStore := "store-456"
		dto := &transfert.Caisse{ID: &idCaisse, StoreID: &idStore}
		storeDTO := &transfert.Store{ID: &idStore}
		caisse := &entities.Caisse{ID: "caisse-123"}

		mockPerms.On("IsGrantedByRoles", []security.Role{user.ROLE_EMPLOYEE}).Return(true)
		mockRepo.On("ReadStore", storeDTO, mock.Anything).Return(&entities.Store{ID: "store-456"}, nil)
		mockRepo.On("CreateCaisse", dto, mock.Anything).Return(caisse, nil)

		result, err := service.CreateCaisse(dto)
		assert.Nil(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, caisse, result)

		mockRepo.AssertExpectations(t)
		mockPerms.AssertExpectations(t)
	})

	t.Run("Devrait retourner une erreur lorsque non autorisé", func(t *testing.T) {
		service, _, mockPerms := setup()

		idCaisse := "caisse-123"
		dto := &transfert.Caisse{ID: &idCaisse}

		mockPerms.On("IsGrantedByRoles", []security.Role{user.ROLE_EMPLOYEE}).Return(false)

		result, err := service.CreateCaisse(dto)
		assert.Nil(t, result)
		assert.NotNil(t, err)
		assert.Equal(t, errors.ErrUnauthorized, err)

		mockPerms.AssertExpectations(t)
	})

	t.Run("Devrait retourner une erreur lorsque dto est nil", func(t *testing.T) {
		service, _, _ := setup()

		result, err := service.CreateCaisse(nil)
		assert.Nil(t, result)
		assert.NotNil(t, err)
		assert.Equal(t, errors.ErrNoDto, err)
	})

	t.Run("Devrait retourner une erreur lorsque store introuvable", func(t *testing.T) {
		service, mockRepo, mockPerms := setup()

		idCaisse := "caisse-123"
		idStore := "store-456"
		dto := &transfert.Caisse{ID: &idCaisse, StoreID: &idStore}
		storeDTO := &transfert.Store{ID: &idStore}

		mockPerms.On("IsGrantedByRoles", []security.Role{user.ROLE_EMPLOYEE}).Return(true)
		mockRepo.On("ReadStore", storeDTO, mock.Anything).Return(nil, errors.ErrNoData)

		result, err := service.CreateCaisse(dto)
		assert.Nil(t, result)
		assert.NotNil(t, err)
		assert.Equal(t, errors.ErrNoData, err)

		mockRepo.AssertExpectations(t)
		mockPerms.AssertExpectations(t)
	})

	t.Run("Devrait retourner une erreur lorsque CreateCaisse échoue", func(t *testing.T) {
		service, mockRepo, mockPerms := setup()

		idCaisse := "caisse-123"
		idStore := "store-456"
		dto := &transfert.Caisse{ID: &idCaisse, StoreID: &idStore}
		storeDTO := &transfert.Store{ID: &idStore}

		mockPerms.On("IsGrantedByRoles", []security.Role{user.ROLE_EMPLOYEE}).Return(true)
		mockRepo.On("ReadStore", storeDTO, mock.Anything).Return(&entities.Store{ID: "store-456"}, nil)
		mockRepo.On("CreateCaisse", dto, mock.Anything).Return(nil, errors.ErrNoData)

		result, err := service.CreateCaisse(dto)
		assert.Nil(t, result)
		assert.NotNil(t, err)
		assert.Equal(t, errors.ErrNoData, err)

		mockRepo.AssertExpectations(t)
		mockPerms.AssertExpectations(t)
	})
}

// Test_UpdateCaisse tests the UpdateCaisse method of StoreService
// Parameters:
// - t: *testing.T
//
// Returns:
// - None: no return value
func Test_UpdateCaisse(t *testing.T) {
	t.Run("Devrait mettre à jour une caisse lorsque autorisé et existe", func(t *testing.T) {
		service, mockRepo, mockPerms := setup()

		idCaisse := "caisse-123"
		dto := &transfert.Caisse{ID: &idCaisse}
		caisse := &entities.Caisse{ID: "caisse-123"}

		mockPerms.On("IsGrantedByRoles", []security.Role{user.ROLE_EMPLOYEE}).Return(true)
		mockRepo.On("ReadCaisse", &transfert.Caisse{ID: &idCaisse}, mock.Anything).Return(caisse, nil)
		mockRepo.On("UpdateCaisse", caisse, mock.Anything).Return(nil)

		result, err := service.UpdateCaisse(dto)
		assert.Nil(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, caisse, result)

		mockRepo.AssertExpectations(t)
		mockPerms.AssertExpectations(t)
	})

	t.Run("Devrait retourner une erreur lorsque non autorisé", func(t *testing.T) {
		service, _, mockPerms := setup()

		idCaisse := "caisse-123"
		dto := &transfert.Caisse{ID: &idCaisse}

		mockPerms.On("IsGrantedByRoles", []security.Role{user.ROLE_EMPLOYEE}).Return(false)

		result, err := service.UpdateCaisse(dto)
		assert.Nil(t, result)
		assert.NotNil(t, err)
		assert.Equal(t, errors.ErrUnauthorized, err)

		mockPerms.AssertExpectations(t)
	})

	t.Run("Devrait retourner une erreur lorsque dto est nil", func(t *testing.T) {
		service, _, _ := setup()

		result, err := service.UpdateCaisse(nil)
		assert.Nil(t, result)
		assert.NotNil(t, err)
		assert.Equal(t, errors.ErrNoDto, err)
	})

	t.Run("Devrait retourner une erreur lorsque caisse introuvable", func(t *testing.T) {
		service, mockRepo, mockPerms := setup()

		idCaisse := "caisse-123"
		dto := &transfert.Caisse{ID: &idCaisse}

		mockPerms.On("IsGrantedByRoles", []security.Role{user.ROLE_EMPLOYEE}).Return(true)
		mockRepo.On("ReadCaisse", &transfert.Caisse{ID: &idCaisse}, mock.Anything).Return(nil, errors.ErrNoData)

		result, err := service.UpdateCaisse(dto)
		assert.Nil(t, result)
		assert.NotNil(t, err)
		assert.Equal(t, errors.ErrNoData, err)

		mockRepo.AssertExpectations(t)
		mockPerms.AssertExpectations(t)
	})

	t.Run("Devrait retourner une erreur lorsque UpdateCaisse échoue", func(t *testing.T) {
		service, mockRepo, mockPerms := setup()

		idCaisse := "caisse-123"
		dto := &transfert.Caisse{ID: &idCaisse}
		caisse := &entities.Caisse{ID: "caisse-123"}

		mockPerms.On("IsGrantedByRoles", []security.Role{user.ROLE_EMPLOYEE}).Return(true)
		mockRepo.On("ReadCaisse", &transfert.Caisse{ID: &idCaisse}, mock.Anything).Return(caisse, nil)
		mockRepo.On("UpdateCaisse", caisse, mock.Anything).Return(errors.ErrNoData)

		result, err := service.UpdateCaisse(dto)
		assert.Nil(t, result)
		assert.NotNil(t, err)
		assert.Equal(t, errors.ErrNoData, err)

		mockRepo.AssertExpectations(t)
		mockPerms.AssertExpectations(t)
	})
}

// Test_DeleteCaisse tests the DeleteCaisse method of StoreService
// Parameters:
// - t: *testing.T
//
// Returns:
// - None: no return value
func Test_DeleteCaisse(t *testing.T) {
	t.Run("Devrait supprimer une caisse lorsque autorisé", func(t *testing.T) {
		service, mockRepo, mockPerms := setup()

		idCaisse := "caisse-123"
		dto := &transfert.Caisse{ID: &idCaisse}

		mockPerms.On("IsGrantedByRoles", []security.Role{user.ROLE_EMPLOYEE}).Return(true)
		mockRepo.On("DeleteCaisse", dto, mock.Anything).Return(nil)

		err := service.DeleteCaisse(dto)
		assert.Nil(t, err)

		mockRepo.AssertExpectations(t)
		mockPerms.AssertExpectations(t)
	})

	t.Run("Devrait retourner une erreur lorsque non autorisé", func(t *testing.T) {
		service, _, mockPerms := setup()

		idCaisse := "caisse-123"
		dto := &transfert.Caisse{ID: &idCaisse}

		mockPerms.On("IsGrantedByRoles", []security.Role{user.ROLE_EMPLOYEE}).Return(false)

		err := service.DeleteCaisse(dto)
		assert.NotNil(t, err)
		assert.Equal(t, errors.ErrUnauthorized, err)

		mockPerms.AssertExpectations(t)
	})

	t.Run("Devrait retourner une erreur lorsque dto est nil", func(t *testing.T) {
		service, _, _ := setup()

		err := service.DeleteCaisse(nil)
		assert.NotNil(t, err)
		assert.Equal(t, errors.ErrNoDto, err)
	})

	t.Run("Devrait retourner une erreur lorsque DeleteCaisse échoue", func(t *testing.T) {
		service, mockRepo, mockPerms := setup()

		idCaisse := "caisse-123"
		dto := &transfert.Caisse{ID: &idCaisse}

		mockPerms.On("IsGrantedByRoles", []security.Role{user.ROLE_EMPLOYEE}).Return(true)
		mockRepo.On("DeleteCaisse", dto, mock.Anything).Return(errors.ErrNoData)

		err := service.DeleteCaisse(dto)
		assert.NotNil(t, err)
		assert.Equal(t, errors.ErrNoData, err)

		mockRepo.AssertExpectations(t)
		mockPerms.AssertExpectations(t)
	})
}
