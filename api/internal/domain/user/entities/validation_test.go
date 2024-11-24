package entities_test

import (
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/google/uuid"
	"github.com/kodmain/thetiptop/api/config"
	transfert "github.com/kodmain/thetiptop/api/internal/application/transfert/user"
	"github.com/kodmain/thetiptop/api/internal/domain/user/entities"
	"github.com/stretchr/testify/assert"
)

func TestNewValidation(t *testing.T) {
	config.Load(aws.String("../../../../config.test.yml"))

	val := entities.CreateValidation(&transfert.Validation{
		Token:    nil,
		ClientID: aws.String("1"),
	})

	assert.Nil(t, val.Token)
	err := val.BeforeCreate(nil)
	assert.NoError(t, err)
	assert.NotNil(t, val.ID)

	err = val.BeforeSave(nil)
	assert.NoError(t, err)

	err = val.BeforeUpdate(nil)
	assert.NoError(t, err)

	val = entities.CreateValidation(&transfert.Validation{
		Token:    aws.String("123456"),
		ClientID: aws.String("1"),
	})

	assert.NotNil(t, val.Token)
	assert.Equal(t, val.IsPublic(), false)
	assert.Empty(t, val.GetOwnerID())
	val.CredentialID = aws.String(uuid.New().String())
	assert.Equal(t, val.GetOwnerID(), *val.CredentialID)
}

func TestNewValidationWithoutCID(t *testing.T) {
	config.Load(aws.String("../../../../config.test.yml"))

	val := entities.CreateValidation(&transfert.Validation{
		Token:    nil,
		ClientID: nil,
	})

	assert.Nil(t, val.Token)
	err := val.BeforeCreate(nil)
	assert.Error(t, err)
	assert.Empty(t, val.ID)

	err = val.BeforeSave(nil)
	assert.NoError(t, err)

	err = val.BeforeUpdate(nil)
	assert.Error(t, err)

	val = entities.CreateValidation(&transfert.Validation{
		Token:    aws.String("123456"),
		ClientID: nil,
	})

	assert.NotNil(t, val.Token)

}
