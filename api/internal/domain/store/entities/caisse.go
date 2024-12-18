package entities

import (
	"time"

	"github.com/google/uuid"
	transfert "github.com/kodmain/thetiptop/api/internal/application/transfert/crm"
	"gorm.io/gorm"
)

type Caisses []*Caisse

type Caisse struct {
	// Gorm model
	ID        string          `gorm:"type:varchar(36);primaryKey;" json:"id"`
	CreatedAt time.Time       `json:"-"`
	UpdatedAt time.Time       `json:"-"`
	DeletedAt *gorm.DeletedAt `gorm:"index" json:"-"`

	// Relations
	StoreID *string `gorm:"type:varchar(36);index;"`
}

func (caisse *Caisse) BeforeCreate(tx *gorm.DB) error {
	id, err := uuid.NewRandom()
	if err != nil {
		return err
	}

	caisse.ID = id.String()

	return nil
}

func CreateCaisse(obj *transfert.Caisse) *Caisse {
	c := &Caisse{
		StoreID: obj.StoreID,
	}

	if obj.ID != nil {
		c.ID = *obj.ID
	}

	return c
}
