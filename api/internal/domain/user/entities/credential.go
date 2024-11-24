package entities

import (
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/google/uuid"
	transfert "github.com/kodmain/thetiptop/api/internal/application/transfert/user"
	"github.com/kodmain/thetiptop/api/internal/infrastructure/security/hash"
	"gorm.io/gorm"
)

type Credential struct {
	ID        string          `gorm:"type:varchar(36);primaryKey;" json:"id"`
	CreatedAt time.Time       `json:"-"`
	UpdatedAt time.Time       `json:"-"`
	DeletedAt *gorm.DeletedAt `gorm:"index" json:"-"`

	Email    *string `gorm:"type:varchar(320);uniqueIndex" json:"email"`
	Password *string `gorm:"type:varchar(255)" json:"-"` // private field
}

func (cred *Credential) CompareHash(password string) bool {
	return hash.CompareHash(cred.Password, aws.String(*cred.Email+":"+password), hash.BCRYPT) == nil
}

func (cred *Credential) BeforeUpdate(tx *gorm.DB) error {
	cred.UpdatedAt = time.Now()
	return nil
}

func (cred *Credential) BeforeCreate(tx *gorm.DB) error {
	id, err := uuid.NewRandom()
	if err != nil {
		return err
	}

	cred.ID = id.String()
	return nil
}

func (cred *Credential) IsPublic() bool {
	return false
}

func (cred *Credential) GetOwnerID() string {
	return cred.ID
}

func CreateCredential(obj *transfert.Credential) *Credential {
	return &Credential{
		Email:    obj.Email,
		Password: obj.Password,
	}
}
