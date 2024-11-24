package game_test

import (
	"net/http"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/kodmain/thetiptop/api/internal/application/services/game"
	transfert "github.com/kodmain/thetiptop/api/internal/application/transfert/game"
	"github.com/kodmain/thetiptop/api/internal/domain/game/entities"
	"github.com/kodmain/thetiptop/api/internal/infrastructure/errors"
	"github.com/stretchr/testify/assert"
)

func TestGetTickets(t *testing.T) {
	t.Run("should return tickets successfully", func(t *testing.T) {
		mockService := new(DomainGameService)
		expectedTickets := []*entities.Ticket{}
		mockService.On("GetTickets").Return(expectedTickets, nil)

		statusCode, response := game.GetTickets(mockService)

		assert.Equal(t, fiber.StatusOK, statusCode)
		assert.Equal(t, expectedTickets, response)
		mockService.AssertCalled(t, "GetTickets")
	})

	t.Run("should return error when service fails", func(t *testing.T) {
		mockService := new(DomainGameService)
		expectedError := errors.ErrBadRequest
		mockService.On("GetTickets").Return(nil, expectedError)

		statusCode, response := game.GetTickets(mockService)

		assert.Equal(t, http.StatusBadRequest, statusCode)
		assert.Error(t, response.(*errors.Error))
		mockService.AssertCalled(t, "GetTickets")
	})
}

func TestGetRandomTicket(t *testing.T) {
	t.Run("should return ticket successfully", func(t *testing.T) {
		mockService := new(DomainGameService)
		expectedTicket := &entities.Ticket{}
		mockService.On("GetRandomTicket").Return(expectedTicket, nil)

		statusCode, response := game.GetRandomTicket(mockService)

		assert.Equal(t, fiber.StatusOK, statusCode)
		assert.Equal(t, expectedTicket, response)
		mockService.AssertCalled(t, "GetRandomTicket")
	})

	t.Run("should return error when service fails", func(t *testing.T) {
		mockService := new(DomainGameService)
		expectedError := errors.ErrBadRequest
		mockService.On("GetRandomTicket").Return(nil, expectedError)

		statusCode, response := game.GetRandomTicket(mockService)

		assert.Equal(t, http.StatusBadRequest, statusCode)
		assert.Error(t, response.(*errors.Error))
		mockService.AssertCalled(t, "GetRandomTicket")
	})
}

func TestUpdateTicket(t *testing.T) {
	t.Run("should update ticket successfully", func(t *testing.T) {
		// Create a mock service
		mockService := new(DomainGameService)
		dtoTicket := &transfert.Ticket{}
		updatedTicket := &entities.Ticket{ID: "1", Token: "updated-token"}

		// Configure the mock to return the updated ticket
		mockService.On("UpdateTicket", dtoTicket).Return(updatedTicket, nil)

		// Call the function under test
		statusCode, response := game.UpdateTicket(mockService, dtoTicket)

		// Assert the results
		assert.Equal(t, fiber.StatusOK, statusCode)
		assert.Equal(t, updatedTicket, response)
		mockService.AssertCalled(t, "UpdateTicket", dtoTicket)
	})

	t.Run("should return error when update fails", func(t *testing.T) {
		// Create a mock service
		mockService := new(DomainGameService)
		dtoTicket := &transfert.Ticket{}
		expectedError := errors.ErrBadRequest

		// Configure the mock to return an error
		mockService.On("UpdateTicket", dtoTicket).Return(nil, expectedError)

		// Call the function under test
		statusCode, response := game.UpdateTicket(mockService, dtoTicket)

		// Assert the results
		assert.Equal(t, http.StatusBadRequest, statusCode)
		assert.Equal(t, expectedError, response)
		mockService.AssertCalled(t, "UpdateTicket", dtoTicket)
	})
}
