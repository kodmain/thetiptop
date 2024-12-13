package entities_test

import (
	"testing"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/google/uuid"
	transfert "github.com/kodmain/thetiptop/api/internal/application/transfert/crm"
	"github.com/kodmain/thetiptop/api/internal/domain/store/entities"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func TestCreateStore(t *testing.T) {
	label := "Test Store"
	isOnline := true

	obj := &transfert.Store{
		Label:    &label,
		IsOnline: &isOnline,
	}

	store := entities.CreateStore(obj)

	assert.NotNil(t, store)
	assert.Equal(t, label, *store.Label)
	assert.Equal(t, isOnline, *store.IsOnline)
}

func TestStore_BeforeCreateAndUpdate(t *testing.T) {
	store := &entities.Store{
		Label:    aws.String("Test Store"),
		IsOnline: aws.Bool(true),
		Caisses: entities.Caisses{
			&entities.Caisse{
				ID: uuid.NewString(),
			},
		},
	}

	err := store.BeforeCreate(nil)
	assert.Nil(t, err)
	assert.NotEmpty(t, store.ID)

	oldUpdatedAt := store.UpdatedAt
	time.Sleep(100 * time.Millisecond)

	err = store.BeforeUpdate(nil)
	assert.Nil(t, err)
	assert.True(t, store.UpdatedAt.After(oldUpdatedAt))
}

func TestStore_AfterFind(t *testing.T) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	assert.Nil(t, err)

	err = db.AutoMigrate(&entities.Store{}, &entities.Caisse{})
	assert.Nil(t, err)

	store := &entities.Store{
		ID:       uuid.NewString(),
		Label:    aws.String("Test Store"),
		IsOnline: aws.Bool(true),
	}

	err = db.Create(&store).Error
	assert.Nil(t, err)

	err = db.First(&store, "id = ?", store.ID).Error
	assert.Nil(t, err)

	assert.NotNil(t, store.Caisses)
	assert.Equal(t, 0, len(store.Caisses))
}

func TestStore_BeforeCreate_WithNilCaisses(t *testing.T) {
	store := &entities.Store{
		Label:    aws.String("Test Store"),
		IsOnline: aws.Bool(true),
		Caisses:  nil,
	}

	err := store.BeforeCreate(nil)
	assert.Nil(t, err)
	assert.NotEmpty(t, store.ID)
}

func TestStore_BeforeUpdate_WithNilCaisses(t *testing.T) {
	store := &entities.Store{
		Label:    aws.String("Test Store"),
		IsOnline: aws.Bool(true),
		Caisses:  nil,
	}

	err := store.BeforeUpdate(nil)
	assert.Nil(t, err)
	assert.NotNil(t, store.UpdatedAt)
}

func TestStore_AfterFind_Error(t *testing.T) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	assert.Nil(t, err)

	// Close the database to simulate an error
	sqlDB, err := db.DB()
	assert.Nil(t, err)
	err = sqlDB.Close()
	assert.Nil(t, err)

	store := &entities.Store{
		ID: uuid.NewString(),
	}

	err = store.AfterFind(db)
	assert.NotNil(t, err)
}

func TestCreateStore_WithID(t *testing.T) {
	id := uuid.NewString()
	label := "Test Store"
	isOnline := true

	obj := &transfert.Store{
		ID:       &id,
		Label:    &label,
		IsOnline: &isOnline,
	}

	store := entities.CreateStore(obj)

	assert.NotNil(t, store)
	assert.Equal(t, id, store.ID)
	assert.Equal(t, label, *store.Label)
	assert.Equal(t, isOnline, *store.IsOnline)
}

func TestCreateStore_WithoutID(t *testing.T) {
	label := "Test Store"
	isOnline := true

	obj := &transfert.Store{
		Label:    &label,
		IsOnline: &isOnline,
	}

	store := entities.CreateStore(obj)

	assert.NotNil(t, store)
	assert.Empty(t, store.ID)
	assert.Equal(t, label, *store.Label)
	assert.Equal(t, isOnline, *store.IsOnline)
}

func TestStore_IsPublic(t *testing.T) {
	store := &entities.Store{}
	assert.False(t, store.IsPublic())
}

func TestStore_GetOwnerID(t *testing.T) {
	store := &entities.Store{}
	assert.Equal(t, "", store.GetOwnerID())
}
