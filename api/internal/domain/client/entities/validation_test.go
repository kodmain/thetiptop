package entities_test

import (
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/kodmain/thetiptop/api/config"
	"github.com/kodmain/thetiptop/api/internal/application/transfert"
	"github.com/kodmain/thetiptop/api/internal/domain/client/entities"
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

}
