package password_test

import (
	"testing"

	"github.com/kodmain/thetiptop/api/internal/infrastructure/security/password"
	"github.com/stretchr/testify/assert"
)

func TestGeneratePassword(t *testing.T) {
	t.Run("Generate password with lowercase characters", func(t *testing.T) {
		password, err := password.GeneratePassword(8, password.Lowercase)
		assert.NoError(t, err)
		assert.Equal(t, 8, len(password))
		assert.Regexp(t, "^[a-z]+$", password)
	})

	t.Run("Generate password with uppercase characters", func(t *testing.T) {
		password, err := password.GeneratePassword(10, password.Uppercase)
		assert.NoError(t, err)
		assert.Equal(t, 10, len(password))
		assert.Regexp(t, "^[A-Z]+$", password)
	})

	t.Run("Generate password with digits", func(t *testing.T) {
		password, err := password.GeneratePassword(12, password.Digits)
		assert.NoError(t, err)
		assert.Equal(t, 12, len(password))
		assert.Regexp(t, "^[0-9]+$", password)
	})

	t.Run("Generate password with special characters", func(t *testing.T) {
		password, err := password.GeneratePassword(15, password.SpecialChars)
		assert.NoError(t, err)
		assert.Equal(t, 15, len(password))
		assert.Regexp(t, "^[!@#$%^&*()]+$", password)
	})

	t.Run("Generate password with all character types", func(t *testing.T) {
		password, err := password.GeneratePassword(20, password.All)
		assert.NoError(t, err)
		assert.Equal(t, 20, len(password))
		assert.Regexp(t, "^[a-zA-Z0-9!@#$%^&*()]+$", password)
	})
}
