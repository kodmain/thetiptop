package entities

import (
	"time"

	"github.com/google/uuid"
	transfert "github.com/kodmain/thetiptop/api/internal/application/transfert/crm"
	"gorm.io/gorm"
)

type Store struct {
	// Gorm model
	ID        string          `gorm:"type:varchar(36);primaryKey;" json:"id"`
	CreatedAt time.Time       `json:"-"`
	UpdatedAt time.Time       `json:"-"`
	DeletedAt *gorm.DeletedAt `gorm:"index" json:"-"`

	Label    *string `gorm:"type:varchar(255);uniqueIndex" json:"label"`
	IsOnline *bool   `gorm:"type:boolean" json:"is_online"`

	Caisses Caisses `gorm:"foreignKey:StoreID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"caisses"`
}

func CreateStore(obj *transfert.Store) *Store {
	t := &Store{
		Label:    obj.Label,
		IsOnline: obj.IsOnline,
	}

	if obj.ID != nil {
		t.ID = *obj.ID
	}

	return t
}

func (store *Store) IsPublic() bool {
	return false
}

func (store *Store) GetOwnerID() string {
	return ""
}

func (store *Store) BeforeUpdate(tx *gorm.DB) error {
	store.UpdatedAt = time.Now()

	for _, caisse := range store.Caisses {
		caisse.StoreID = &store.ID
	}

	return nil
}

func (store *Store) BeforeCreate(tx *gorm.DB) error {
	id, err := uuid.NewRandom()
	if err != nil {
		return err
	}

	store.ID = id.String()

	for _, caisse := range store.Caisses {
		caisse.StoreID = &store.ID
	}

	return nil
}

func (store *Store) AfterFind(tx *gorm.DB) error {
	if err := tx.Model(store).Association("Caisses").Find(&store.Caisses); err != nil {
		return err
	}
	return nil
}
