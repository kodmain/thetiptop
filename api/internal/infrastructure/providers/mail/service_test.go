package mail_test

import (
	"testing"

	"github.com/kodmain/thetiptop/api/internal/infrastructure/providers/mail"
	"github.com/stretchr/testify/assert"
)

func TestService(t *testing.T) {

	err := mail.New(nil)
	assert.Error(t, err)

	err = mail.New(&mail.Service{})

	assert.Error(t, err)

	err = mail.New(&mail.Service{
		Host:      "smtp.world.com",
		Port:      "0",
		Username:  "secret",
		Password:  "secret",
		From:      "hello@world.com",
		Expeditor: "Whoami",
		Disable:   true,
	})

	assert.Nil(t, err)

}
