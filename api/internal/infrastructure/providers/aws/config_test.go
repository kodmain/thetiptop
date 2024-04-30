package aws_test

import (
	"testing"

	"github.com/kodmain/thetiptop/api/internal/infrastructure/providers/aws"

	"github.com/stretchr/testify/assert"
)

func TestConnect(t *testing.T) {
	// Call the Connect function
	client, err := aws.Connect("hello")
	assert.Error(t, err)
	assert.Nil(t, client)

	client, err = aws.Connect()
	assert.NoError(t, err)
	assert.NotNil(t, client)
}
