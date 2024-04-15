package jwt_test

import (
	"testing"

	"github.com/kodmain/thetiptop/api/internal/infrastructure/serializers/jwt"
	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {

	err := jwt.New(nil)
	assert.NoError(t, err)

	jwtConfig := &jwt.JWT{
		TZ:       "",
		Secret:   "",
		Expire:   0,
		Refresh:  0,
		Duration: 0,
	}

	err = jwt.New(jwtConfig)
	if err != nil {
		t.Errorf("Failed to create JWT instance: %v", err)
	}

}

func TestFromID(t *testing.T) {
	err := jwt.New(nil)
	assert.NoError(t, err)

	id := "exampleID"

	token, err := jwt.FromID(id)
	assert.NoError(t, err)
	assert.NotEmpty(t, token)

	claims, err := jwt.TokenToClaims(token)
	assert.NoError(t, err)
	assert.Equal(t, id, claims.ID)

	claims, err = jwt.TokenToClaims("fail" + token + "fail")
	assert.Error(t, err)
	assert.Nil(t, claims)
}
