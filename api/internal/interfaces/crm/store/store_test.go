package store_test

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/kodmain/thetiptop/api/internal/domain/store/entities"
	"github.com/kodmain/thetiptop/api/internal/infrastructure/serializers/jwt"
	"github.com/stretchr/testify/assert"
)

func TestStores(t *testing.T) {
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
		encodingName := "FormURLEncoded"
		if encoding == JSONEncoded {
			encodingName = "JSONEncoded"
		}

		t.Run("Store/"+encodingName, func(t *testing.T) {
			content, status, err := request("GET", DOMAIN+"/store", authorization, encoding, nil)
			assert.Nil(t, err)
			assert.Equal(t, http.StatusOK, status)
			var stores []*entities.Store
			assert.Nil(t, json.Unmarshal(content, &stores), "Response should be valid JSON")
			assert.NotNil(t, stores, "Response should not be nil")

			for _, store := range stores {
				content, status, err := request("GET", DOMAIN+"/store/"+store.ID, authorization, encoding, nil)
				assert.Nil(t, err)

				assert.Equal(t, http.StatusOK, status)
				var storeByID *entities.Store
				assert.Nil(t, json.Unmarshal(content, &storeByID), "Response should be valid JSON")
				assert.NotNil(t, storeByID, "Response should not be nil")
			}
		})
	}

	assert.Nil(t, stop())
}
