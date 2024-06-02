package database

import (
	"gorm.io/gorm"
)

type Database struct {
	Config *Config
	Engine *gorm.DB
}
