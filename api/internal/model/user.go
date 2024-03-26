package model

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Account model info
// @Description User account information
// @Description with user id and username
type User struct {
	gorm.Model `swaggerignore:"true"`

	ID         uuid.UUID `gorm:"type:uuid;default:gen_random_uuid()"`
	Email      string    `json:"email" gorm:"unique;not null"`
	Password   string    `gorm:"not null"`
	Newsletter bool      `json:"newsletter" gorm:"not null"`
	CGU        bool      `json:"cgu" gorm:"not null"`

	CreatedAt int64 `json:"createdAt" gorm:"autoCreateTime"`
	UpdatedAt int64 `json:"updatedAt" gorm:"autoUpdateTime"`
}
