package entities

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/kodmain/thetiptop/api/config"
	"github.com/kodmain/thetiptop/api/internal/application/transfert"
	"github.com/kodmain/thetiptop/api/internal/infrastructure/security/token"
	"gorm.io/gorm"
)

type Validation struct {
	// gorm model
	ID        string         `gorm:"type:varchar(36);primaryKey;"`
	CreatedAt time.Time      `json:"-"`
	UpdatedAt time.Time      `json:"-"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
	// entity
	Token     token.Luhn     `gorm:"type:varchar(6);index"`
	Type      ValidationType `gorm:"type:varchar(10)"`
	Validated bool           `gorm:"type:boolean;default:false"`
	ClientID  string         `gorm:"type:varchar(36);uniqueIndex" json:"-"`
	ExpiresAt time.Time      `json:"-"`
}

// BeforeSave attribue la valeur de ClientID avant de sauvegarder
func (v *Validation) BeforeSave(tx *gorm.DB) error {
	if v.Token == "" {
		v.Token = token.Generate(6)
	}

	return nil
}

func (validation *Validation) BeforeCreate(tx *gorm.DB) error {
	id, err := uuid.NewRandom()
	if err != nil {
		return err
	}

	if validation.ClientID == "" {
		return errors.New("ClientID is required on Validation")
	}

	validation.ID = id.String()
	duration, err := time.ParseDuration(config.Get("security.validation.expire").(string))
	if err != nil {
		return err
	}

	validation.ExpiresAt = time.Now().Add(duration)

	return nil
}

func (validation *Validation) BeforeUpdate(tx *gorm.DB) error {
	validation.UpdatedAt = time.Now()
	if validation.ClientID == "" {
		return errors.New("Client is required on Validation")
	}

	return nil
}

func CreateValidation(obj *transfert.Validation) *Validation {
	return &Validation{
		Token: *obj.Token,
	}
}
