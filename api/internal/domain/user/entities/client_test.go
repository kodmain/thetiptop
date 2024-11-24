package entities_test

import (
	"testing"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/google/uuid"
	transfert "github.com/kodmain/thetiptop/api/internal/application/transfert/user"
	"github.com/kodmain/thetiptop/api/internal/domain/user/entities"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func TestClient_HasSuccessValidation(t *testing.T) {
	client := &entities.Client{
		Validations: entities.Validations{
			&entities.Validation{
				Type:      entities.PasswordRecover,
				Validated: true,
			},
			&entities.Validation{
				Type:      entities.PasswordRecover,
				Validated: false,
			},
			&entities.Validation{
				Type:      entities.PasswordRecover,
				Validated: true,
			},
		},
	}

	validationType := entities.PasswordRecover
	result := client.HasSuccessValidation(entities.PasswordRecover)

	assert.NotNil(t, result)
	assert.Equal(t, validationType, result.Type)
	assert.True(t, result.Validated)

	result = client.HasSuccessValidation(entities.MailValidation)
	assert.Nil(t, result)

	assert.Equal(t, client.IsPublic(), false)
	assert.Equal(t, client.GetOwnerID(), client.ID)
	client.CredentialID = aws.String(uuid.New().String())
	assert.Equal(t, client.GetOwnerID(), *client.CredentialID)
}

func TestClient_HasNotExpiredValidation(t *testing.T) {
	client := &entities.Client{
		Validations: entities.Validations{
			&entities.Validation{
				Type:      entities.PasswordRecover,
				Validated: true,
				ExpiresAt: time.Now().Add(time.Hour),
			},
			&entities.Validation{
				Type:      entities.PasswordRecover,
				Validated: false,
				ExpiresAt: time.Now().Add(-time.Hour),
			},
			&entities.Validation{
				Type:      entities.PasswordRecover,
				Validated: false,
				ExpiresAt: time.Now().Add(time.Hour),
			},
		},
	}

	validationType := entities.PasswordRecover
	result := client.HasNotExpiredValidation(entities.PasswordRecover)

	assert.NotNil(t, result)
	assert.Equal(t, validationType, result.Type)
	assert.False(t, result.Validated)
	assert.False(t, result.HasExpired())

	result = client.HasNotExpiredValidation(entities.MailValidation)
	assert.Nil(t, result)
}

func TestClientBeforeCreateAndUpdate(t *testing.T) {
	client := &entities.Client{
		CGU:        aws.Bool(true),
		Newsletter: aws.Bool(false),
	}

	client.Validations = append(client.Validations, &entities.Validation{
		Type: entities.MailValidation,
	})

	// Test BeforeCreate
	err := client.BeforeCreate(nil)
	assert.Nil(t, err)
	assert.NotEmpty(t, client.ID)
	assert.Equal(t, client.ID, client.ID)

	// Simule un timestamp précédent pour comparer avec AfterUpdate
	old := client.UpdatedAt
	time.Sleep(100 * time.Millisecond)

	// Test BeforeUpdate
	err = client.BeforeUpdate(nil)
	assert.Nil(t, err)
	assert.True(t, client.UpdatedAt.After(old))
}

func TestClient_AfterFind(t *testing.T) {
	// Setup SQLite in-memory database for testing
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	assert.Nil(t, err)

	// Automigrate for Client and Validation tables
	err = db.AutoMigrate(&entities.Client{}, &entities.Validation{})
	assert.Nil(t, err)

	// Create a new client with no validations
	client := &entities.Client{
		ID: uuid.New().String(),
	}

	// Insert the client into the database
	err = db.Create(&client).Error
	assert.Nil(t, err)

	// Call AfterFind by retrieving the client from the database
	err = db.First(&client, "id = ?", client.ID).Error
	assert.Nil(t, err)

	// Test if the Validations relation is properly loaded
	assert.NotNil(t, client.Validations)
	assert.True(t, len(client.Validations) == 0) // The client should have 0 validations
}

func TestCreateClient(t *testing.T) {
	// Cas de test avec des valeurs valides pour CGU et Newsletter
	cgu := aws.Bool(true)
	newsletter := aws.Bool(false)

	// Crée un objet de type transfert.Client avec les champs nécessaires
	input := &transfert.Client{
		CGU:        cgu,
		Newsletter: newsletter,
	}

	// Appelle la fonction CreateClient
	client := entities.CreateClient(input)

	// Vérifie que le client n'est pas nil
	assert.NotNil(t, client)

	// Vérifie que les champs CGU et Newsletter ont été correctement transférés
	assert.Equal(t, cgu, client.CGU)
	assert.Equal(t, newsletter, client.Newsletter)

	// Vérifie que le champ Validations est bien initialisé à une liste vide
	assert.NotNil(t, client.Validations)
	assert.Equal(t, 0, len(client.Validations))
}

func TestCreateClientWithNilFields(t *testing.T) {
	// Cas de test où CGU et Newsletter sont nil
	input := &transfert.Client{
		CGU:        nil,
		Newsletter: nil,
	}

	// Appelle la fonction CreateClient
	client := entities.CreateClient(input)

	// Vérifie que le client n'est pas nil
	assert.NotNil(t, client)

	// Vérifie que les champs CGU et Newsletter sont nil
	assert.Nil(t, client.CGU)
	assert.Nil(t, client.Newsletter)

	// Vérifie que le champ Validations est bien initialisé à une liste vide
	assert.NotNil(t, client.Validations)
	assert.Equal(t, 0, len(client.Validations))
}
