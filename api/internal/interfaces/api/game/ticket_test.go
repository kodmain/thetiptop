package game_test

import (
	"encoding/json"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/kodmain/thetiptop/api/internal/domain/game/entities"
	"github.com/kodmain/thetiptop/api/internal/infrastructure/serializers/jwt"
	"github.com/stretchr/testify/assert"
)

func TestGame(t *testing.T) {
	encodingTypes := []EncodingType{FormURLEncoded, JSONEncoded}
	assert.Nil(t, start(8888, 8444))

	JWT, status, err := request("POST", "http://localhost:8888/user/auth", "", JSONEncoded, map[string][]any{
		"email":    {email},
		"password": {password},
	})

	assert.Nil(t, err)
	assert.Equal(t, 200, status)
	assert.NotNil(t, JWT)

	var tokenData fiber.Map
	err = json.Unmarshal(JWT, &tokenData)

	accessTokenString := tokenData["access_token"].(string)
	authorization := "Bearer " + accessTokenString
	assert.Nil(t, err)
	claims, err := jwt.TokenToClaims(accessTokenString)
	assert.Nil(t, err)
	assert.NotNil(t, claims)

	for _, encoding := range encodingTypes {
		var encodingName string = "FormURLEncoded"
		if encoding == JSONEncoded {
			encodingName = "JSONEncoded"
		}

		t.Run("GetTicket/"+encodingName, func(t *testing.T) {
			randomTicket, status, err := request("GET", "http://localhost:8888/game/ticket/random", authorization, encoding)
			assert.Nil(t, err)
			assert.Equal(t, 200, status)

			ticket := entities.Ticket{}
			json.Unmarshal(randomTicket, &ticket)

			assert.NotNil(t, ticket)
			t.Run("UpdateTicket/"+encodingName, func(t *testing.T) {
				updatedTicket, status, err := request("PUT", "http://localhost:8888/game/ticket", authorization, encoding, map[string][]any{
					"id": {ticket.ID},
				})
				assert.Nil(t, err)
				assert.Equal(t, 200, status)

				ticket = entities.Ticket{}
				json.Unmarshal(updatedTicket, &ticket)

				assert.NotNil(t, ticket)

				t.Run("GetTickets/"+encodingName, func(t *testing.T) {
					tickets, status, err := request("GET", "http://localhost:8888/game/ticket", authorization, encoding)
					assert.Nil(t, err)
					assert.Equal(t, 200, status)

					ticketsArray := []*entities.Ticket{}
					json.Unmarshal(tickets, &ticketsArray)

					assert.NotNil(t, ticketsArray)
				})
			})
		})

	}

	assert.Nil(t, stop())
}
