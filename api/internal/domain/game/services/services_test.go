package services_test

import (
	"github.com/kodmain/thetiptop/api/internal/application/security"
	"github.com/kodmain/thetiptop/api/internal/domain/game/services"
	"github.com/kodmain/thetiptop/api/internal/infrastructure/providers/database"
	"github.com/stretchr/testify/mock"
)

// GameRepositoryMock est le mock pour GameRepositoryInterface
type GameRepositoryMock struct {
	mock.Mock
}

// PermissionMock est le mock pour PermissionInterface
type PermissionMock struct {
	mock.Mock
}

func (m *PermissionMock) IsGranted(roles ...string) bool {
	args := m.Called(roles)
	return args.Bool(0)
}

func (m *PermissionMock) CanRead(ressource database.Entity, rules ...security.Rule) bool {
	args := m.Called(ressource)
	return args.Bool(0)
}

func (m *PermissionMock) CanCreate(ressource database.Entity, rules ...security.Rule) bool {
	args := m.Called(ressource)
	return args.Bool(0)
}

func (m *PermissionMock) CanUpdate(ressource database.Entity, rules ...security.Rule) bool {
	args := m.Called(ressource)
	return args.Bool(0)
}

func (m *PermissionMock) CanDelete(ressource database.Entity, rules ...security.Rule) bool {
	args := m.Called(ressource)
	return args.Bool(0)
}

func PermissionLock() {}

func setup() *GameRepositoryMock {
	mockRepository := new(GameRepositoryMock)
	mockSecurity := new(PermissionMock)

	service := services.Game(mockSecurity, mockRepository)

	return service, mockRepository
}
