package entities_test

import (
	"testing"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/google/uuid"
	transfert "github.com/kodmain/thetiptop/api/internal/application/transfert/game"
	"github.com/kodmain/thetiptop/api/internal/domain/game/entities"
	"github.com/kodmain/thetiptop/api/internal/infrastructure/security/token"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func TestCreateTicket(t *testing.T) {
	// Cas avec des données valides
	input := &transfert.Ticket{
		ID:       aws.String(uuid.New().String()),
		ClientID: aws.String(uuid.New().String()),
		Prize:    aws.String("PrizeA"),
		Token:    aws.String("123456"),
	}

	ticket := entities.CreateTicket(input)

	assert.NotNil(t, ticket)
	assert.Equal(t, *input.ID, ticket.ID)
	assert.Equal(t, *input.ClientID, *ticket.ClientID)
	assert.Equal(t, *input.Prize, *ticket.Prize)
	assert.Equal(t, token.Luhn("123456"), ticket.Token)
}

func TestCreateTicketWithNilFields(t *testing.T) {
	// Cas avec des champs optionnels nuls
	input := &transfert.Ticket{
		ID:       nil,
		ClientID: nil,
		Prize:    nil,
		Token:    aws.String("123456"),
	}

	ticket := entities.CreateTicket(input)

	assert.NotNil(t, ticket)
	assert.Empty(t, ticket.ID)
	assert.Nil(t, ticket.ClientID)
	assert.Nil(t, ticket.Prize)
	assert.Equal(t, token.Luhn("123456"), ticket.Token)
}

func TestTicket_IsPublic(t *testing.T) {
	ticket := &entities.Ticket{}
	assert.False(t, ticket.IsPublic())
}

func TestTicket_GetOwnerID(t *testing.T) {
	t.Run("with ClientID", func(t *testing.T) {
		clientID := uuid.New().String()
		ticket := &entities.Ticket{
			ClientID: aws.String(clientID),
		}
		assert.Equal(t, clientID, ticket.GetOwnerID())
	})

	t.Run("without ClientID", func(t *testing.T) {
		ticket := &entities.Ticket{
			ClientID: nil,
		}
		assert.Equal(t, "", ticket.GetOwnerID())
	})
}

func TestTicket_BeforeCreate(t *testing.T) {
	ticket := &entities.Ticket{}
	err := ticket.BeforeCreate(nil)

	assert.Nil(t, err)
	assert.NotEmpty(t, ticket.ID)
}

func TestTicket_BeforeUpdate(t *testing.T) {
	ticket := &entities.Ticket{
		UpdatedAt: time.Now().Add(-time.Hour), // Ancienne date
	}
	oldTime := ticket.UpdatedAt

	err := ticket.BeforeUpdate(nil)
	assert.Nil(t, err)
	assert.True(t, ticket.UpdatedAt.After(oldTime))
}

func TestTicket_AfterFind(t *testing.T) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	assert.Nil(t, err)

	// Migrer la table Ticket
	err = db.AutoMigrate(&entities.Ticket{})
	assert.Nil(t, err)

	// Insérer un ticket dans la base
	ticket := &entities.Ticket{
		ID: uuid.New().String(),
	}
	err = db.Create(ticket).Error
	assert.Nil(t, err)

	// Lire le ticket depuis la base pour déclencher AfterFind
	var fetchedTicket entities.Ticket
	err = db.First(&fetchedTicket, "id = ?", ticket.ID).Error
	assert.Nil(t, err)
	assert.Equal(t, ticket.ID, fetchedTicket.ID)
}
