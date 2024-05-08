package database_test

import (
	"github.com/kodmain/thetiptop/api/internal/infrastructure/providers/database"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

type ServiceMock struct {
	Config *database.Config
	db     mock.Mock
}

func (s *ServiceMock) AutoMigrate(values ...interface{}) error {
	args := s.db.Called(values)
	return args.Error(0)
}

func (s *ServiceMock) Create(value interface{}) *gorm.DB {
	args := s.db.Called(value)
	return args.Get(0).(*gorm.DB)
}

func (s *ServiceMock) Where(query interface{}, args ...interface{}) *gorm.DB {
	args = append([]interface{}{query}, args...)
	call := s.db.Called(args...)
	return call.Get(0).(*gorm.DB)
}
