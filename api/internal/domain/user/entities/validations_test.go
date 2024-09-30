package entities

import (
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestHasInClient(t *testing.T) {
	v := &Validation{
		Type: PasswordRecover,
	}

	vs := Validations{v}

	assert.NotNil(t, vs.Has(PasswordRecover))
	assert.Nil(t, vs.Has(MailValidation))

	assert.Equal(t, v.IsPublic(), false)
	assert.Empty(t, v.GetOwnerID())
	v.CredentialID = aws.String(uuid.New().String())
	assert.Equal(t, v.GetOwnerID(), *v.CredentialID)

}
