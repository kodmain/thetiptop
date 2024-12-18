package entities_test

import (
	"testing"

	"github.com/google/uuid"
	transfert "github.com/kodmain/thetiptop/api/internal/application/transfert/crm"
	"github.com/kodmain/thetiptop/api/internal/domain/store/entities"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func TestCreateCaisse(t *testing.T) {
	storeID := uuid.NewString()
	obj := &transfert.Caisse{
		StoreID: &storeID,
	}

	caisse := entities.CreateCaisse(obj)

	assert.NotNil(t, caisse)
	assert.Equal(t, storeID, *caisse.StoreID)
}

func TestCreateCaisse_WithID(t *testing.T) {
	id := uuid.NewString()
	storeID := uuid.NewString()
	obj := &transfert.Caisse{
		ID:      &id,
		StoreID: &storeID,
	}

	caisse := entities.CreateCaisse(obj)

	assert.NotNil(t, caisse)
	assert.Equal(t, id, caisse.ID)
	assert.Equal(t, storeID, *caisse.StoreID)
}

func TestCaisse_BeforeCreate(t *testing.T) {
	caisse := &entities.Caisse{
		StoreID: nil,
	}

	err := caisse.BeforeCreate(nil)
	assert.Nil(t, err)
	assert.NotEmpty(t, caisse.ID)
}

func TestCaisse_CRUD(t *testing.T) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	assert.Nil(t, err)

	err = db.AutoMigrate(&entities.Caisse{})
	assert.Nil(t, err)

	caisse := &entities.Caisse{
		StoreID: nil,
	}

	err = db.Create(&caisse).Error
	assert.Nil(t, err)
	assert.NotEmpty(t, caisse.ID)

	var fetchedCaisse entities.Caisse
	err = db.First(&fetchedCaisse, "id = ?", caisse.ID).Error
	assert.Nil(t, err)
	assert.Equal(t, caisse.ID, fetchedCaisse.ID)
}
