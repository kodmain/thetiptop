package services_test

import (
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/kodmain/thetiptop/api/internal/application/security"
	transfert "github.com/kodmain/thetiptop/api/internal/application/transfert/game"
	"github.com/kodmain/thetiptop/api/internal/domain/game/entities"
	user "github.com/kodmain/thetiptop/api/internal/domain/user/entities"
	"github.com/kodmain/thetiptop/api/internal/infrastructure/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func Test_GetRandomTicket(t *testing.T) {
	t.Run("Should return ticket", func(t *testing.T) {
		service, mockRepo, mockPerms := setup()

		mockRepo.On("ReadTicket", &transfert.Ticket{}, mock.Anything).Return(&entities.Ticket{}, nil)
		mockPerms.On("IsGrantedByRoles", []security.Role{user.ROLE_EMPLOYEE}).Return(true)

		ticket, err := service.GetRandomTicket()
		assert.Nil(t, err)
		assert.NotNil(t, ticket)

		mockRepo.AssertExpectations(t)
	})

	t.Run("Should return error when unauthorized", func(t *testing.T) {
		service, _, mockPerms := setup()

		mockPerms.On("IsGrantedByRoles", []security.Role{user.ROLE_EMPLOYEE}).Return(false)

		ticket, err := service.GetRandomTicket()
		assert.NotNil(t, err)
		assert.Nil(t, ticket)

		mockPerms.AssertExpectations(t)
	})

	t.Run("Should return error when repository return error", func(t *testing.T) {
		service, mockRepo, mockPerms := setup()

		mockRepo.On("ReadTicket", &transfert.Ticket{}, mock.Anything).Return(nil, errors.ErrNoData)
		mockPerms.On("IsGrantedByRoles", []security.Role{user.ROLE_EMPLOYEE}).Return(true)

		ticket, err := service.GetRandomTicket()
		assert.NotNil(t, err)
		assert.Nil(t, ticket)

		mockRepo.AssertExpectations(t)
	})
}

func Test_GetTickets(t *testing.T) {
	t.Run("Should return tickets", func(t *testing.T) {
		service, mockRepo, mockPerms := setup()

		credentialID := "valid-credential-id"

		// Configuration des mocks
		mockRepo.On("ReadTickets", mock.Anything, mock.Anything).Return([]*entities.Ticket{{ID: "123"}}, nil)
		mockPerms.On("GetCredentialID").Return(&credentialID)

		// Appeler la méthode testée
		tickets, err := service.GetTickets()

		// Assertions
		assert.Nil(t, err)        // Vérifie qu'il n'y a pas d'erreur
		assert.NotNil(t, tickets) // Vérifie que des tickets sont retournés
		assert.Len(t, tickets, 1) // Vérifie le nombre de tickets

		// Vérifications des attentes
		mockRepo.AssertExpectations(t)
		mockPerms.AssertExpectations(t)
	})

	t.Run("Should return error when repository return error", func(t *testing.T) {
		service, mockRepo, mockPerms := setup()

		credentialID := "valid-credential-id"

		// Configuration des mocks
		mockPerms.On("GetCredentialID").Return(&credentialID)
		mockRepo.On("ReadTickets", mock.Anything, mock.Anything).Return(nil, errors.ErrNoData)

		// Appeler la méthode testée
		tickets, err := service.GetTickets()

		// Assertions
		assert.NotNil(t, err)
		assert.Nil(t, tickets)
		assert.Equal(t, errors.ErrNoData, err)

		// Vérifications des attentes
		mockRepo.AssertExpectations(t)
		mockPerms.AssertExpectations(t)
	})
}

func Test_UpdateTicket(t *testing.T) {
	cid := aws.String("client-123")
	t.Run("Should return updated ticket", func(t *testing.T) {
		service, mockRepo, mockPerms := setup()

		dto := &transfert.Ticket{
			CredentialID: cid,
		}

		ticket := &entities.Ticket{
			ID:           "ticket-123",
			CredentialID: cid,
		}

		mockRepo.On("ReadTicket", dto, mock.Anything).Return(ticket, nil)
		mockPerms.On("IsAuthenticated").Return(true)
		mockPerms.On("GetCredentialID").Return(cid)
		mockRepo.On("UpdateTicket", ticket, mock.Anything).Return(nil)

		updatedTicket, err := service.UpdateTicket(dto)
		assert.Nil(t, err)
		assert.NotNil(t, updatedTicket)

		assert.Equal(t, ticket, updatedTicket)

		mockRepo.AssertExpectations(t)
		mockPerms.AssertExpectations(t)
	})

	t.Run("Should return error when ticket not found", func(t *testing.T) {
		service, mockRepo, mockPerms := setup()

		dto := &transfert.Ticket{
			CredentialID: cid,
		}

		mockRepo.On("ReadTicket", dto, mock.Anything).Return(nil, errors.ErrNoData)
		mockPerms.On("IsAuthenticated").Return(true)

		updatedTicket, err := service.UpdateTicket(dto)
		assert.NotNil(t, err)
		assert.Nil(t, updatedTicket)
		assert.Equal(t, errors.ErrNoData, err)

		mockRepo.AssertExpectations(t)
	})

	t.Run("Should return error when unauthorized", func(t *testing.T) {
		service, mockRepo, mockPerms := setup()

		dto := &transfert.Ticket{
			CredentialID: cid,
		}

		ticket := &entities.Ticket{
			ID:           "ticket-123",
			CredentialID: cid,
		}

		mockRepo.On("ReadTicket", dto, mock.Anything).Return(ticket, nil)
		mockPerms.On("IsAuthenticated").Return(false)

		updatedTicket, err := service.UpdateTicket(dto)
		assert.NotNil(t, err)
		assert.Nil(t, updatedTicket)
		assert.Equal(t, errors.ErrUnauthorized, err)

		mockRepo.AssertExpectations(t)
		mockPerms.AssertExpectations(t)
	})

	t.Run("Should return error when update fails", func(t *testing.T) {
		service, mockRepo, mockPerms := setup()

		dto := &transfert.Ticket{
			CredentialID: cid,
		}

		ticket := &entities.Ticket{
			ID:           "ticket-123",
			CredentialID: cid,
		}

		// Configuration des mocks
		mockRepo.On("ReadTicket", dto, mock.Anything).Return(ticket, nil)
		mockPerms.On("IsAuthenticated").Return(true)
		mockPerms.On("GetCredentialID").Return(cid)
		mockRepo.On("UpdateTicket", ticket, mock.Anything).Return(errors.ErrNoData)

		// Appel de la méthode à tester
		updatedTicket, err := service.UpdateTicket(dto)

		// Vérification des résultats
		assert.NotNil(t, err)
		assert.Nil(t, updatedTicket)
		assert.Equal(t, errors.ErrNoData, err)

		mockRepo.AssertExpectations(t)
		mockPerms.AssertExpectations(t)
	})

}

func Test_GetTicketById(t *testing.T) {
	t.Run("Should return ticket when authorized and ticket exists", func(t *testing.T) {
		service, mockRepo, mockPerms := setup()

		dto := &transfert.Ticket{
			ID: aws.String("ticket-123"),
		}

		ticket := &entities.Ticket{
			ID: "ticket-123",
		}

		// Configuration des mocks
		mockPerms.On("IsGrantedByRoles", []security.Role{user.ROLE_EMPLOYEE}).Return(true)
		mockRepo.On("ReadTicket", dto, mock.Anything).Return(ticket, nil)

		// Appel de la méthode à tester
		result, err := service.GetTicketById(dto)

		// Vérifications des résultats
		assert.Nil(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, ticket, result)

		mockRepo.AssertExpectations(t)
		mockPerms.AssertExpectations(t)
	})

	t.Run("Should return error when unauthorized", func(t *testing.T) {
		service, _, mockPerms := setup()

		dto := &transfert.Ticket{
			ID: aws.String("ticket-123"),
		}

		// Configuration des mocks
		mockPerms.On("IsGrantedByRoles", []security.Role{user.ROLE_EMPLOYEE}).Return(false)

		// Appel de la méthode à tester
		result, err := service.GetTicketById(dto)

		// Vérifications des résultats
		assert.NotNil(t, err)
		assert.Nil(t, result)
		assert.Equal(t, errors.ErrUnauthorized, err)

		mockPerms.AssertExpectations(t)
	})

	t.Run("Should return error when ticket not found", func(t *testing.T) {
		service, mockRepo, mockPerms := setup()

		dto := &transfert.Ticket{
			ID: aws.String("ticket-123"),
		}

		// Configuration des mocks
		mockPerms.On("IsGrantedByRoles", []security.Role{user.ROLE_EMPLOYEE}).Return(true)
		mockRepo.On("ReadTicket", dto, mock.Anything).Return(nil, errors.ErrNoData)

		// Appel de la méthode à tester
		result, err := service.GetTicketById(dto)

		// Vérifications des résultats
		assert.NotNil(t, err)
		assert.Nil(t, result)
		assert.Equal(t, errors.ErrNoData, err)

		mockRepo.AssertExpectations(t)
		mockPerms.AssertExpectations(t)
	})

	t.Run("Should return error when repository fails", func(t *testing.T) {
		service, mockRepo, mockPerms := setup()

		dto := &transfert.Ticket{
			ID: aws.String("ticket-123"),
		}

		// Configuration des mocks
		mockPerms.On("IsGrantedByRoles", []security.Role{user.ROLE_EMPLOYEE}).Return(true)
		mockRepo.On("ReadTicket", dto, mock.Anything).Return(nil, errors.ErrBadRequest)

		// Appel de la méthode à tester
		result, err := service.GetTicketById(dto)

		// Vérifications des résultats
		assert.NotNil(t, err)
		assert.Nil(t, result)
		assert.Equal(t, errors.ErrBadRequest, err)

		mockRepo.AssertExpectations(t)
		mockPerms.AssertExpectations(t)
	})

}
