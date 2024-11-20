package services_test

import (
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	transfert "github.com/kodmain/thetiptop/api/internal/application/transfert/game"
	"github.com/kodmain/thetiptop/api/internal/domain/game/entities"
	"github.com/kodmain/thetiptop/api/internal/infrastructure/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func Test_GetRandomTicket(t *testing.T) {
	t.Run("Should return ticket", func(t *testing.T) {
		service, repo, _ := setup()

		repo.On("ReadTicket", &transfert.Ticket{}, mock.Anything).Return(&entities.Ticket{}, nil)

		ticket, err := service.GetRandomTicket()
		assert.Nil(t, err)
		assert.NotNil(t, ticket)

		repo.AssertExpectations(t)
	})

	t.Run("Should return error when repository return error", func(t *testing.T) {
		service, repo, _ := setup()

		repo.On("ReadTicket", &transfert.Ticket{}, mock.Anything).Return(nil, errors.ErrNoData)

		ticket, err := service.GetRandomTicket()
		assert.NotNil(t, err)
		assert.Nil(t, ticket)

		repo.AssertExpectations(t)
	})
}

func Test_GetTickets(t *testing.T) {
	t.Run("Should return tickets", func(t *testing.T) {
		service, repo, permission := setup()

		credentialID := "valid-credential-id"

		permission.On("GetCredentialID").Return(&credentialID)
		repo.On("ReadTickets", mock.Anything, mock.Anything).Return([]*entities.Ticket{{ID: "123"}}, nil)

		// Appeler la méthode testée
		tickets, err := service.GetTickets()
		assert.Nil(t, err)        // Vérifie qu'il n'y a pas d'erreur
		assert.NotNil(t, tickets) // Vérifie que des tickets sont retournés
		assert.Len(t, tickets, 1) // Vérifie le nombre de tickets

		// Vérifications des attentes
		repo.AssertExpectations(t)
		permission.AssertExpectations(t)
	})

	t.Run("Should return error when repository return error", func(t *testing.T) {
		service, repo, permission := setup()

		credentialID := "valid-credential-id"

		permission.On("GetCredentialID").Return(&credentialID)
		repo.On("ReadTickets", mock.Anything, mock.Anything).Return(nil, errors.ErrNoData)

		tickets, err := service.GetTickets()
		assert.NotNil(t, err)
		assert.Nil(t, tickets)
		assert.Equal(t, errors.ErrNoData, err)

		repo.AssertExpectations(t)
		permission.AssertExpectations(t)
	})
}

func Test_UpdateTicket(t *testing.T) {
	cid := aws.String("client-123")
	t.Run("Should update ticket successfully", func(t *testing.T) {
		service, repo, sec := setup()

		dto := &transfert.Ticket{
			ClientID: cid,
		}

		ticket := &entities.Ticket{
			ID:       "ticket-123",
			ClientID: cid,
		}

		// Configuration correcte des appels mock
		repo.On("ReadTicket", dto, mock.Anything).Return(ticket, nil)
		sec.On("CanUpdate", ticket, mock.Anything).Return(true)
		repo.On("UpdateTicket", ticket, mock.Anything).Return(nil)

		// Appel de la méthode à tester
		updatedTicket, err := service.UpdateTicket(dto)

		// Assertions pour valider le comportement
		assert.Nil(t, err)
		assert.NotNil(t, updatedTicket)
		assert.Equal(t, dto.ClientID, updatedTicket.ClientID)

		repo.AssertExpectations(t)
		sec.AssertExpectations(t)
	})

	t.Run("Should return error when ticket not found", func(t *testing.T) {
		service, repo, _ := setup()

		dto := &transfert.Ticket{
			ClientID: cid,
		}

		repo.On("ReadTicket", dto, mock.Anything).Return(nil, errors.ErrNoData)

		updatedTicket, err := service.UpdateTicket(dto)
		assert.NotNil(t, err)
		assert.Nil(t, updatedTicket)
		assert.Equal(t, errors.ErrNoData, err)

		repo.AssertExpectations(t)
	})

	t.Run("Should return error when unauthorized", func(t *testing.T) {
		service, repo, sec := setup()

		dto := &transfert.Ticket{
			ClientID: cid,
		}

		ticket := &entities.Ticket{
			ID:       "ticket-123",
			ClientID: cid,
		}

		repo.On("ReadTicket", dto, mock.Anything).Return(ticket, nil)
		sec.On("CanUpdate", ticket, mock.Anything).Return(false)

		updatedTicket, err := service.UpdateTicket(dto)
		assert.NotNil(t, err)
		assert.Nil(t, updatedTicket)
		assert.Equal(t, errors.ErrUnauthorized, err)

		repo.AssertExpectations(t)
		sec.AssertExpectations(t)
	})

	t.Run("Should return error when update fails", func(t *testing.T) {
		service, repo, sec := setup()

		dto := &transfert.Ticket{
			ClientID: cid,
		}

		ticket := &entities.Ticket{
			ID:       "ticket-123",
			ClientID: cid,
		}

		// Configuration des mocks
		repo.On("ReadTicket", dto, mock.Anything).Return(ticket, nil)
		sec.On("CanUpdate", ticket, mock.Anything).Return(true)
		repo.On("UpdateTicket", ticket, mock.Anything).Return(errors.ErrNoData)

		// Appel de la méthode à tester
		updatedTicket, err := service.UpdateTicket(dto)

		// Vérification des résultats
		assert.NotNil(t, err)
		assert.Nil(t, updatedTicket)
		assert.Equal(t, errors.ErrNoData, err)

		repo.AssertExpectations(t)
		sec.AssertExpectations(t)
	})

}
