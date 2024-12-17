package code_test

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestListErrors(t *testing.T) {
	encodingTypes := []EncodingType{FormURLEncoded, JSONEncoded}
	assert.Nil(t, start(8888, 8444))

	for _, encoding := range encodingTypes {
		var encodingName string = "FormURLEncoded"
		if encoding == JSONEncoded {
			encodingName = "JSONEncoded"
		}

		t.Run("ListErrors/"+encodingName, func(t *testing.T) {
			// Effectuer la requête GET /code/error
			response, status, err := request("GET", "http://localhost:8888/code/error", "", encoding)
			assert.Nil(t, err)
			assert.Equal(t, http.StatusOK, status)

			// Parse de la réponse en map[string]interface{} pour vérifier les contenus
			var errorMap map[string]interface{}
			err = json.Unmarshal(response, &errorMap)
			assert.Nil(t, err)
			assert.NotNil(t, errorMap)
		})
	}

	assert.Nil(t, stop())
}
