package database

import (
	"gorm.io/gorm"
)

type ServiceInterface interface {
	AutoMigrate(...interface{}) error
	Create(interface{}) *gorm.DB
	Where(query interface{}, args ...interface{}) *gorm.DB
}

type Service struct {
	Config *Config
	db     *gorm.DB
}

func (s *Service) AutoMigrate(values ...interface{}) error {
	return s.db.AutoMigrate(values...)
}

func (s *Service) Create(value interface{}) *gorm.DB {
	return s.db.Create(value)
}

func (s *Service) Where(query interface{}, args ...interface{}) *gorm.DB {
	return s.db.Where(query, args...)
}
