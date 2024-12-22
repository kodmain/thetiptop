// services_test.go
//
// Description:
// Ce fichier contient les mocks pour le StoreService, incluant un mock pour StoreRepositoryInterface
// et pour PermissionInterface. Il fournit également une fonction setup() permettant de créer
// une instance du service, ainsi que des mocks de permission et de repository.
//
// Les commentaires dans le code sont en anglais et respectent le format demandé.
//
// Parameters:
// - None
//
// Returns:
// - None: no return value

package services_test

import (
	"github.com/kodmain/thetiptop/api/internal/application/security"
	transfert "github.com/kodmain/thetiptop/api/internal/application/transfert/crm"
	"github.com/kodmain/thetiptop/api/internal/domain/store/entities"
	"github.com/kodmain/thetiptop/api/internal/domain/store/services"
	"github.com/kodmain/thetiptop/api/internal/infrastructure/errors"
	"github.com/kodmain/thetiptop/api/internal/infrastructure/providers/database"
	"github.com/stretchr/testify/mock"
)

// StoreRepositoryMock is the mock for StoreRepositoryInterface
//
// Parameters:
// - None
//
// Returns:
// - None: no return value
type StoreRepositoryMock struct {
	mock.Mock
}

// CreateStores simulates creating multiple stores in the repository
// Parameters:
// - objs: []*transfert.Store, the store dtos to create
// - options: ...database.Option, additional database options
//
// Returns:
// - errors.ErrorInterface: an error if creation fails
func (m *StoreRepositoryMock) CreateStores(objs []*transfert.Store, options ...database.Option) errors.ErrorInterface {
	// CreateStores(obj []*transfert.Store, options ...database.Option)
	args := m.Called(objs, options)
	if args.Get(0) != nil {
		return args.Error(0).(errors.ErrorInterface)
	}
	return nil
}

// ReadStores simulates reading multiple stores from the repository
// Parameters:
// - obj: *transfert.Store, the store dto to read
// - options: ...database.Option, additional database options
//
// Returns:
// - []*entities.Store: the store entities if found
// - errors.ErrorInterface: an error if reading fails
func (m *StoreRepositoryMock) ReadStores(obj *transfert.Store, options ...database.Option) ([]*entities.Store, errors.ErrorInterface) {
	args := m.Called(obj, options)
	if args.Get(0) == nil {
		return nil, args.Error(1).(errors.ErrorInterface)
	}
	return args.Get(0).([]*entities.Store), nil
}

// ReadStore simulates reading a store from the repository
// Parameters:
// - dto: *transfert.Store, the store dto to read
// - options: ...database.Option, additional database options
//
// Returns:
// - *entities.Store: the store entity if found
// - errors.ErrorInterface: an error if reading fails
func (m *StoreRepositoryMock) ReadStore(dto *transfert.Store, options ...database.Option) (*entities.Store, errors.ErrorInterface) {
	args := m.Called(dto, options)
	if args.Get(0) == nil {
		return nil, args.Error(1).(errors.ErrorInterface)
	}
	return args.Get(0).(*entities.Store), nil
}

// DeleteStores simulates deleting multiple stores in the repository
// Parameters:
// - objs: []*transfert.Store, the store dtos to delete
// - options: ...database.Option, additional database options
//
// Returns:
// - errors.ErrorInterface: an error if deletion fails
func (m *StoreRepositoryMock) DeleteStores(objs []*transfert.Store, options ...database.Option) errors.ErrorInterface {
	args := m.Called(objs, options)
	if args.Get(0) != nil {
		return args.Error(0).(errors.ErrorInterface)
	}
	return nil
}

// UpdateStores simulates updating multiple stores in the repository
// Parameters:
// - objs: []*entities.Store, the store entities to update
// - options: ...database.Option, additional database options
//
// Returns:
// - errors.ErrorInterface: an error if update fails
func (m *StoreRepositoryMock) UpdateStores(objs []*entities.Store, options ...database.Option) errors.ErrorInterface {
	args := m.Called(objs, options)
	if args.Get(0) != nil {
		return args.Error(0).(errors.ErrorInterface)
	}
	return nil
}

// CreateCaisse simulates creating a caisse in the repository
// Parameters:
// - obj: *transfert.Caisse, the caisse dto to create
// - options: ...database.Option, additional database options
//
// Returns:
// - *entities.Caisse: the created caisse entity
// - errors.ErrorInterface: an error if creation fails
func (m *StoreRepositoryMock) CreateCaisse(obj *transfert.Caisse, options ...database.Option) (*entities.Caisse, errors.ErrorInterface) {
	args := m.Called(obj, options)
	if args.Get(0) == nil {
		return nil, args.Error(1).(errors.ErrorInterface)
	}
	return args.Get(0).(*entities.Caisse), nil
}

// ReadCaisse simulates reading a caisse from the repository
// Parameters:
// - obj: *transfert.Caisse, the caisse dto to read
// - options: ...database.Option, additional database options
//
// Returns:
// - *entities.Caisse: the caisse entity if found
// - errors.ErrorInterface: an error if reading fails
func (m *StoreRepositoryMock) ReadCaisse(obj *transfert.Caisse, options ...database.Option) (*entities.Caisse, errors.ErrorInterface) {
	args := m.Called(obj, options)
	if args.Get(0) == nil {
		return nil, args.Error(1).(errors.ErrorInterface)
	}
	return args.Get(0).(*entities.Caisse), nil
}

// ReadCaisses simulates reading multiple caisses from the repository
// Parameters:
// - obj: *transfert.Caisse, the caisse dto to read
// - options: ...database.Option, additional database options
//
// Returns:
// - []*entities.Caisse: the caisse entities if found
// - errors.ErrorInterface: an error if reading fails
func (m *StoreRepositoryMock) ReadCaisses(obj *transfert.Caisse, options ...database.Option) ([]*entities.Caisse, errors.ErrorInterface) {
	args := m.Called(obj, options)
	if args.Get(0) == nil {
		return nil, args.Error(1).(errors.ErrorInterface)
	}
	return args.Get(0).([]*entities.Caisse), nil
}

// DeleteCaisse simulates deleting a caisse in the repository
// Parameters:
// - obj: *transfert.Caisse, the caisse dto to delete
// - options: ...database.Option, additional database options
//
// Returns:
// - errors.ErrorInterface: an error if deletion fails
func (m *StoreRepositoryMock) DeleteCaisse(obj *transfert.Caisse, options ...database.Option) errors.ErrorInterface {
	args := m.Called(obj, options)
	if args.Get(0) != nil {
		return args.Error(0).(errors.ErrorInterface)
	}
	return nil
}

// UpdateCaisse simulates updating a caisse in the repository
// Parameters:
// - obj: *entities.Caisse, the caisse entity to update
// - options: ...database.Option, additional database options
//
// Returns:
// - errors.ErrorInterface: an error if update fails
func (m *StoreRepositoryMock) UpdateCaisse(obj *entities.Caisse, options ...database.Option) errors.ErrorInterface {
	args := m.Called(obj, options)
	if args.Get(0) != nil {
		return args.Error(0).(errors.ErrorInterface)
	}
	return nil
}

// PermissionMock is the mock for PermissionInterface
//
// Parameters:
// - None
//
// Returns:
// - None: no return value
type PermissionMock struct {
	mock.Mock
}

// IsAuthenticated simulates checking if a user is authenticated
// Returns:
// - bool: true if authenticated, false otherwise
func (m *PermissionMock) IsAuthenticated() bool {
	args := m.Called()
	return args.Bool(0)
}

// IsGrantedByRoles simulates checking if a user has required roles
// Parameters:
// - roles: ...security.Role, roles required
//
// Returns:
// - bool: true if roles are granted, false otherwise
func (m *PermissionMock) IsGrantedByRoles(roles ...security.Role) bool {
	args := m.Called(roles)
	return args.Bool(0)
}

// IsGrantedByRules simulates checking if a user has required rules
// Parameters:
// - rules: ...security.Rule, rules required
//
// Returns:
// - bool: true if rules are granted, false otherwise
func (m *PermissionMock) IsGrantedByRules(rules ...security.Rule) bool {
	args := m.Called(rules)
	return args.Bool(0)
}

// CanRead simulates checking if a user can read a given resource
// Parameters:
// - resource: database.Entity, the resource to check
//
// Returns:
// - bool: true if user can read, false otherwise
func (m *PermissionMock) CanRead(resource database.Entity, rules ...security.Rule) bool {
	args := m.Called(resource)
	return args.Bool(0)
}

// CanCreate simulates checking if a user can create a given resource
// Parameters:
// - resource: database.Entity, the resource to check
//
// Returns:
// - bool: true if user can create, false otherwise
func (m *PermissionMock) CanCreate(resource database.Entity, rules ...security.Rule) bool {
	args := m.Called(resource)
	return args.Bool(0)
}

// CanUpdate simulates checking if a user can update a given resource
// Parameters:
// - resource: database.Entity, the resource to check
//
// Returns:
// - bool: true if user can update, false otherwise
func (m *PermissionMock) CanUpdate(resource database.Entity, rules ...security.Rule) bool {
	args := m.Called(resource)
	return args.Bool(0)
}

// CanDelete simulates checking if a user can delete a given resource
// Parameters:
// - resource: database.Entity, the resource to check
//
// Returns:
// - bool: true if user can delete, false otherwise
func (m *PermissionMock) CanDelete(resource database.Entity, rules ...security.Rule) bool {
	args := m.Called(resource)
	return args.Bool(0)
}

// GetCredentialID simulates retrieving the credential ID of the current user
// Returns:
// - *string: the credential ID if available
func (m *PermissionMock) GetCredentialID() *string {
	args := m.Called()
	if args.Get(0) == nil {
		return nil
	}
	return args.Get(0).(*string)
}

// setup function initializes a StoreService with mocked repository and permissions
// Parameters:
// - None
//
// Returns:
// - *services.StoreService: an instance of StoreService
// - *StoreRepositoryMock: a mock of the store repository
// - *PermissionMock: a mock of the permission interface
func setup() (*services.StoreService, *StoreRepositoryMock, *PermissionMock) {
	mockRepository := new(StoreRepositoryMock)
	mockSecurity := new(PermissionMock)

	service := services.Store(mockSecurity, mockRepository)

	return service, mockRepository, mockSecurity
}
