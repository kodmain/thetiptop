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

// Test_ListStores tests the ListStores method of StoreService
// Parameters:
// - t: *testing.T
//
// Returns:
// - None: no return value
func Test_ListStores(t *testing.T) {
	t.Run("Devrait retourner les stores lorsque autorisé et disponible", func(t *testing.T) {
		service, mockRepo, mockPerms := setup()

		stores := []*entities.Store{
			{ID: "store-1"},
			{ID: "store-2"},
		}

		mockPerms.On("IsGrantedByRoles", []security.Role{user.ROLE_EMPLOYEE}).Return(true)
		mockRepo.On("ReadStores", &transfert.Store{}, mock.Anything).Return(stores, nil)

		result, err := service.ListStores()
		assert.Nil(t, err)
		assert.NotNil(t, result)
		assert.Len(t, result, 2)

		mockRepo.AssertExpectations(t)
		mockPerms.AssertExpectations(t)
	})

	t.Run("Devrait retourner une erreur lorsque non autorisé", func(t *testing.T) {
		service, _, mockPerms := setup()

		mockPerms.On("IsGrantedByRoles", []security.Role{user.ROLE_EMPLOYEE}).Return(false)

		result, err := service.ListStores()
		assert.Nil(t, result)
		assert.NotNil(t, err)
		assert.Equal(t, errors.ErrUnauthorized, err)

		mockPerms.AssertExpectations(t)
	})

	t.Run("Devrait retourner une erreur lorsque le repo retourne une erreur", func(t *testing.T) {
		service, mockRepo, mockPerms := setup()

		mockPerms.On("IsGrantedByRoles", []security.Role{user.ROLE_EMPLOYEE}).Return(true)
		mockRepo.On("ReadStores", &transfert.Store{}, mock.Anything).Return(nil, errors.ErrNoData)

		result, err := service.ListStores()
		assert.Nil(t, result)
		assert.NotNil(t, err)
		assert.Equal(t, errors.ErrNoData, err)

		mockRepo.AssertExpectations(t)
		mockPerms.AssertExpectations(t)
	})
}

// Test_GetStoreByID tests the GetStoreByID method of StoreService
// Parameters:
// - t: *testing.T
//
// Returns:
// - None: no return value
func Test_GetStoreByID(t *testing.T) {
	t.Run("Devrait retourner un store lorsque autorisé et existe", func(t *testing.T) {
		service, mockRepo, mockPerms := setup()

		idStore := "store-123"
		dto := &transfert.Store{ID: &idStore}
		store := &entities.Store{ID: "store-123"}

		mockPerms.On("IsGrantedByRoles", []security.Role{user.ROLE_EMPLOYEE}).Return(true)
		mockRepo.On("ReadStore", dto, mock.Anything).Return(store, nil)

		result, err := service.GetStoreByID(dto)
		assert.Nil(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, store, result)

		mockRepo.AssertExpectations(t)
		mockPerms.AssertExpectations(t)
	})

	t.Run("Devrait retourner une erreur lorsque non autorisé", func(t *testing.T) {
		service, _, mockPerms := setup()

		idStore := "store-123"
		dto := &transfert.Store{ID: &idStore}

		mockPerms.On("IsGrantedByRoles", []security.Role{user.ROLE_EMPLOYEE}).Return(false)

		result, err := service.GetStoreByID(dto)
		assert.Nil(t, result)
		assert.NotNil(t, err)
		assert.Equal(t, errors.ErrUnauthorized, err)

		mockPerms.AssertExpectations(t)
	})

	t.Run("Devrait retourner une erreur lorsque dto est nil", func(t *testing.T) {
		service, _, _ := setup()

		result, err := service.GetStoreByID(nil)
		assert.Nil(t, result)
		assert.NotNil(t, err)
		assert.Equal(t, errors.ErrNoDto, err)
	})

	t.Run("Devrait retourner une erreur lorsque le repo échoue", func(t *testing.T) {
		service, mockRepo, mockPerms := setup()

		idStore := "store-123"
		dto := &transfert.Store{ID: &idStore}

		mockPerms.On("IsGrantedByRoles", []security.Role{user.ROLE_EMPLOYEE}).Return(true)
		mockRepo.On("ReadStore", dto, mock.Anything).Return(nil, errors.ErrNoData)

		result, err := service.GetStoreByID(dto)
		assert.Nil(t, result)
		assert.NotNil(t, err)
		assert.Equal(t, errors.ErrNoData, err)

		mockRepo.AssertExpectations(t)
		mockPerms.AssertExpectations(t)
	})
}
