package events_test

import (
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	transfert "github.com/kodmain/thetiptop/api/internal/application/transfert/crm"
	"github.com/kodmain/thetiptop/api/internal/domain/store/entities"
	"github.com/kodmain/thetiptop/api/internal/domain/store/events"
	"github.com/kodmain/thetiptop/api/internal/infrastructure/errors"
	"github.com/kodmain/thetiptop/api/internal/infrastructure/providers/database"
)

// MockStoreRepository est une implémentation mock de StoreRepositoryInterface utilisant testify/mock
type MockStoreRepository struct {
	mock.Mock
}

// ReadStores simule la méthode ReadStores de StoreRepositoryInterface
func (m *MockStoreRepository) ReadStores(obj *transfert.Store, options ...database.Option) ([]*entities.Store, errors.ErrorInterface) {
	args := m.Called(obj, options)

	var err errors.ErrorInterface
	if e := args.Get(1); e != nil {
		err = e.(errors.ErrorInterface)
	}

	if result := args.Get(0); result != nil {
		return result.([]*entities.Store), err
	}

	return nil, err
}

// CreateCaisse simule la méthode CreateCaisse de StoreRepositoryInterface
func (m *MockStoreRepository) CreateCaisse(obj *transfert.Caisse, options ...database.Option) (*entities.Caisse, errors.ErrorInterface) {
	args := m.Called(obj, options)

	var err errors.ErrorInterface
	if e := args.Get(1); e != nil {
		err = e.(errors.ErrorInterface)
	}

	if result := args.Get(0); result != nil {
		return result.(*entities.Caisse), err
	}

	return nil, err
}

// ReadCaisses simule la méthode ReadCaisses de StoreRepositoryInterface
func (m *MockStoreRepository) ReadCaisses(obj *transfert.Caisse, options ...database.Option) ([]*entities.Caisse, errors.ErrorInterface) {
	args := m.Called(obj, options)

	var err errors.ErrorInterface
	if e := args.Get(1); e != nil {
		err = e.(errors.ErrorInterface)
	}

	if result := args.Get(0); result != nil {
		return result.([]*entities.Caisse), err
	}

	return nil, err
}

// ReadCaisse simule la méthode ReadCaisse de StoreRepositoryInterface
func (m *MockStoreRepository) ReadCaisse(obj *transfert.Caisse, options ...database.Option) (*entities.Caisse, errors.ErrorInterface) {
	args := m.Called(obj, options)

	var err errors.ErrorInterface
	if e := args.Get(1); e != nil {
		err = e.(errors.ErrorInterface)
	}

	if result := args.Get(0); result != nil {
		return result.(*entities.Caisse), err
	}

	return nil, err
}

// UpdateCaisse simule la méthode UpdateCaisse de StoreRepositoryInterface
func (m *MockStoreRepository) UpdateCaisse(obj *entities.Caisse, options ...database.Option) errors.ErrorInterface {
	args := m.Called(obj, options)

	if args.Error(0) != nil {
		return args.Error(0).(errors.ErrorInterface)
	}

	return nil
}

// DeleteCaisse simule la méthode DeleteCaisse de StoreRepositoryInterface
func (m *MockStoreRepository) DeleteCaisse(obj *transfert.Caisse, options ...database.Option) errors.ErrorInterface {
	args := m.Called(obj, options)

	if args.Error(0) != nil {
		return args.Error(0).(errors.ErrorInterface)
	}

	return nil
}

// CreateStores simule la méthode CreateStores de StoreRepositoryInterface
func (m *MockStoreRepository) CreateStores(objs []*transfert.Store, options ...database.Option) errors.ErrorInterface {
	args := m.Called(objs, options)

	if args.Error(0) != nil {
		return args.Error(0).(errors.ErrorInterface)
	}

	return nil
}

// DeleteStores simule la méthode DeleteStores de StoreRepositoryInterface
func (m *MockStoreRepository) DeleteStores(objs []*transfert.Store, options ...database.Option) errors.ErrorInterface {
	args := m.Called(objs, options)

	if args.Error(0) != nil {
		return args.Error(0).(errors.ErrorInterface)
	}

	return nil
}

// UpdateStores simule la méthode UpdateStores de StoreRepositoryInterface
func (m *MockStoreRepository) UpdateStores(objs []*entities.Store, options ...database.Option) errors.ErrorInterface {
	args := m.Called(objs, options)

	if args.Error(0) != nil {
		return args.Error(0).(errors.ErrorInterface)
	}

	return nil
}

// ReadStore simule la méthode ReadStore de StoreRepositoryInterface
func (m *MockStoreRepository) ReadStore(obj *transfert.Store, options ...database.Option) (*entities.Store, errors.ErrorInterface) {
	args := m.Called(obj, options)

	var err errors.ErrorInterface
	if e := args.Get(1); e != nil {
		err = e.(errors.ErrorInterface)
	}

	if result := args.Get(0); result != nil {
		return result.(*entities.Store), err
	}

	return nil, err
}

// setupStoreRepository initialise l'environnement de test en créant une instance du mock et en retournant une fonction de nettoyage
func setupStoreRepository() (*MockStoreRepository, func()) {
	mockRepo := new(MockStoreRepository)
	cleanup := func() {}
	return mockRepo, cleanup
}

// TestConvertEntityToTransfer teste la fonction ConvertEntityToTransfer
func TestConvertEntityToTransfer(t *testing.T) {
	t.Run("convert entity to transfer", func(t *testing.T) {
		// Création d'une entité Store
		entityStore := &entities.Store{
			Label:    aws.String("Store One"),
			IsOnline: aws.Bool(true),
		}

		// Appel de la fonction de conversion
		transferStore := events.ConvertEntityToTransfer(entityStore)

		// Assertions
		assert.NotNil(t, transferStore)
		assert.Equal(t, entityStore.Label, transferStore.Label)
		assert.Equal(t, entityStore.IsOnline, transferStore.IsOnline)
	})
}

// TestConvertTransferToEntity teste la fonction ConvertTransferToEntity
func TestConvertTransferToEntity(t *testing.T) {
	t.Run("convert transfer to entity", func(t *testing.T) {
		// Création d'un transfert Store
		transferStore := &transfert.Store{
			Label:    aws.String("Store Two"),
			IsOnline: aws.Bool(false),
		}

		// Appel de la fonction de conversion
		entityStore := events.ConvertTransferToEntity(transferStore)

		// Assertions
		assert.NotNil(t, entityStore)
		assert.Equal(t, transferStore.Label, entityStore.Label)
		assert.Equal(t, transferStore.IsOnline, entityStore.IsOnline)
	})
}

func TestCreateStores(t *testing.T) {

	t.Run("nominal case", func(t *testing.T) {
		mockRepo, cleanup := setupStoreRepository()
		defer cleanup()

		mockRepo.On("ReadStores", mock.Anything, mock.Anything).Return([]*entities.Store{}, nil)
		mockRepo.On("CreateStores", mock.Anything, mock.Anything).Return(nil)
		mockRepo.On("UpdateStores", mock.Anything, mock.Anything).Return(nil)
		events.CreateStores(mockRepo)
	})

	t.Run("nominal existing store", func(t *testing.T) {
		mockRepo, cleanup := setupStoreRepository()
		defer cleanup()

		mockRepo.On("ReadStores", mock.Anything, mock.Anything).Return([]*entities.Store{
			{
				Label:    aws.String("DigitalStore"),
				IsOnline: aws.Bool(true),
			},
			{
				Label:    aws.String("NotDesiredStore"),
				IsOnline: aws.Bool(false),
			},
		}, nil)
		mockRepo.On("CreateStores", mock.Anything, mock.Anything).Return(nil)
		mockRepo.On("UpdateStores", mock.Anything, mock.Anything).Return(nil)
		mockRepo.On("DeleteStores", mock.Anything, mock.Anything).Return(nil)
		events.CreateStores(mockRepo)
	})

	t.Run("error ReadStores", func(t *testing.T) {
		mockRepo, cleanup := setupStoreRepository()
		defer cleanup()

		// Configuration du mock pour ReadStores avec une erreur
		mockRepo.On("ReadStores", mock.Anything, mock.Anything).Return(nil, errors.ErrInternalServer).Once()

		// Configuration des autres mocks (non appelés dans ce scénario)
		mockRepo.On("CreateStores", mock.Anything, mock.Anything).Return(nil)
		mockRepo.On("UpdateStores", mock.Anything, mock.Anything).Return(nil)

		// Assertion que la fonction panique avec le message attendu
		assert.PanicsWithValue(t, "Failed to read stores: common.internal_error", func() {
			events.CreateStores(mockRepo)
		})

		// Assertions pour vérifier que les autres méthodes ne sont pas appelées
		mockRepo.AssertCalled(t, "ReadStores", mock.Anything, mock.Anything)
		mockRepo.AssertNotCalled(t, "CreateStores", mock.Anything, mock.Anything)
		mockRepo.AssertNotCalled(t, "UpdateStores", mock.Anything, mock.Anything)
		mockRepo.AssertNotCalled(t, "DeleteStores", mock.Anything, mock.Anything)
	})
}
