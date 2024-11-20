package game_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGame(t *testing.T) {
	encodingTypes := []EncodingType{FormURLEncoded, JSONEncoded}
	assert.Nil(t, start(8888, 8444))

	for _, encoding := range encodingTypes {
		var encodingName string = "FormURLEncoded"
		if encoding == JSONEncoded {
			encodingName = "JSONEncoded"
		}

		t.Run("SignUp/"+encodingName, func(t *testing.T) {

		})

		request("DELETE", "http://0.0.0.0:1080/email/all", "", encoding)

	}
}
