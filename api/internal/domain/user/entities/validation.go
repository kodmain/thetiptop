package entities

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/kodmain/thetiptop/api/config"
	transfert "github.com/kodmain/thetiptop/api/internal/application/transfert/user"
	"github.com/kodmain/thetiptop/api/internal/infrastructure/security/token"
	"gorm.io/gorm"
)

type Validation struct {
	// gorm model
	ID        string         `gorm:"type:varchar(36);primaryKey;" json:"id"`
	CreatedAt time.Time      `json:"-"`
	UpdatedAt time.Time      `json:"-"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
	// entity
	Token     *token.Luhn    `gorm:"type:varchar(6);index" json:"token"`
	Type      ValidationType `gorm:"type:varchar(10)" json:"type"`
	Validated bool           `gorm:"type:boolean;default:false" json:"validated"`

	ClientID   *string `gorm:"type:varchar(36)" json:"-"`
	EmployeeID *string `gorm:"type:varchar(36)" json:"-"`

	CredentialID *string   `gorm:"type:varchar(36);index;" json:"-"` // Foreign key to Credential
	ExpiresAt    time.Time `json:"-"`
}

func (v *Validation) HasExpired() bool {
	if v.ExpiresAt.IsZero() {
		return false
	}

	return v.ExpiresAt.Before(time.Now())
}

// BeforeSave attribue la valeur de ClientID avant de sauvegarder
func (v *Validation) BeforeSave(tx *gorm.DB) error {
	if v.Token == nil {
		v.Token = token.Generate(6).Pointer()
	}

	return nil
}

func (validation *Validation) BeforeCreate(tx *gorm.DB) error {
	id, err := uuid.NewRandom()
	if err != nil {
		return err
	}

	if validation.ClientID == nil && validation.EmployeeID == nil {
		return fmt.Errorf("ClientID or EmployeeID is required on Validation")
	}

	validation.ID = id.String()
	duration, err := time.ParseDuration(config.Get("security.validation.expire", nil).(string))
	if err != nil {
		return err
	}

	validation.ExpiresAt = time.Now().Add(duration)

	return nil
}

func (validation *Validation) BeforeUpdate(tx *gorm.DB) error {
	validation.UpdatedAt = time.Now()
	if validation.ClientID == nil && validation.EmployeeID == nil {
		return fmt.Errorf("ClientID or EmployeeID is required on Validation")
	}

	return nil
}

func (validation *Validation) IsPublic() bool {
	return false
}

func (validation *Validation) GetOwnerID() string {
	if validation.CredentialID == nil {
		return ""
	}

	return *validation.CredentialID
}

func CreateValidation(obj *transfert.Validation) *Validation {
	v := &Validation{
		ClientID:   obj.ClientID,
		EmployeeID: obj.EmployeeID,
	}

	if obj.Type != nil {
		if vt, err := newValidationType(obj.Type); err == nil {
			v.Type = vt
		}
	}

	if obj.Token != nil {
		v.Token = token.NewLuhn(*obj.Token).Pointer()
	}

	return v
}
