package services_test

import (
	transfert "github.com/kodmain/thetiptop/api/internal/application/transfert/crm"
	"github.com/kodmain/thetiptop/api/internal/domain/store/entities"
	"github.com/kodmain/thetiptop/api/internal/infrastructure/errors"
	"github.com/stretchr/testify/mock"
)

// MockStoreService est une implémentation mock de StoreServiceInterface utilisant testify/mock
type MockStoreService struct {
	mock.Mock
}

// ListStores simule la méthode ListStores de StoreServiceInterface
func (m *MockStoreService) ListStores() ([]*entities.Store, errors.ErrorInterface) {
	args := m.Called()
	if result := args.Get(0); result != nil {
		return result.([]*entities.Store), nil
	}
	return nil, args.Get(1).(errors.ErrorInterface)
}

// GetStoreByID simule la méthode GetStoreByID de StoreServiceInterface
func (m *MockStoreService) GetStoreByID(dtoStore *transfert.Store) (*entities.Store, errors.ErrorInterface) {
	args := m.Called(dtoStore)
	if result := args.Get(0); result != nil {
		return result.(*entities.Store), nil
	}
	return nil, args.Get(1).(errors.ErrorInterface)
}

// GetCaisse simule la méthode GetCaisse de StoreServiceInterface
func (m *MockStoreService) GetCaisse(dtoCaisse *transfert.Caisse) (*entities.Caisse, errors.ErrorInterface) {
	args := m.Called(dtoCaisse)
	if result := args.Get(0); result != nil {
		return result.(*entities.Caisse), nil
	}
	return nil, args.Get(1).(errors.ErrorInterface)
}

// CreateCaisse simule la méthode CreateCaisse de StoreServiceInterface
func (m *MockStoreService) CreateCaisse(dtoCaisse *transfert.Caisse) (*entities.Caisse, errors.ErrorInterface) {
	args := m.Called(dtoCaisse)
	if result := args.Get(0); result != nil {
		return result.(*entities.Caisse), nil
	}
	return nil, args.Get(1).(errors.ErrorInterface)
}

// DeleteCaisse simule la méthode DeleteCaisse de StoreServiceInterface
func (m *MockStoreService) DeleteCaisse(dtoCaisse *transfert.Caisse) errors.ErrorInterface {
	args := m.Called(dtoCaisse)
	if args.Get(0) == nil {
		return nil
	}

	return args.Get(0).(errors.ErrorInterface)
}

// UpdateCaisse simule la méthode UpdateCaisse de StoreServiceInterface
func (m *MockStoreService) UpdateCaisse(dtoCaisse *transfert.Caisse) (*entities.Caisse, errors.ErrorInterface) {
	args := m.Called(dtoCaisse)
	if result := args.Get(0); result != nil {
		return result.(*entities.Caisse), nil
	}
	return nil, args.Get(1).(errors.ErrorInterface)
}

// setup initialise l'environnement de test en créant une instance du mock et en retournant une fonction de nettoyage
//
// Parameters:
// - None
//
// Returns:
// - *MockStoreService: instance du service mocké
// - func(): fonction de nettoyage (noop dans ce cas)
func setup() (*MockStoreService, func()) {
	mockService := new(MockStoreService)
	cleanup := func() {}
	return mockService, cleanup
}
