package jwt_test

import (
	"testing"

	"github.com/kodmain/thetiptop/api/internal/infrastructure/serializers/jwt"
	"github.com/stretchr/testify/assert"
)

func TestData(t *testing.T) {
	t.Run("with args", func(t *testing.T) {
		jwt.New(nil)

		token, refresh, err := jwt.FromID("oki", nil)
		assert.NoError(t, err)

		access, err := jwt.TokenToClaims(token)
		assert.NoError(t, err)
		assert.NotNil(t, access)
		assert.NotNil(t, refresh)
	})

	t.Run("without args", func(t *testing.T) {
		jwt.New(nil)

		token, refresh, err := jwt.FromID("oki", map[string]any{
			"type": "oki",
		})
		assert.NoError(t, err)

		access, err := jwt.TokenToClaims(token)
		assert.NoError(t, err)
		assert.NotNil(t, access)
		assert.NotNil(t, refresh)
	})
}
