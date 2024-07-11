package entities

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHasInClient(t *testing.T) {
	v := &Validation{
		Type: PasswordRecover,
	}

	vs := Validations{v}

	assert.NotNil(t, vs.Has(PasswordRecover))
	assert.Nil(t, vs.Has(MailValidation))

}
