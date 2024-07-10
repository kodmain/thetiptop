package mail_test

import (
	"testing"

	"github.com/kodmain/thetiptop/api/internal/infrastructure/providers/mail"
	"github.com/stretchr/testify/assert"
)

func TestService(t *testing.T) {
	service := &mail.Service{}
	m := &mail.Mail{}

	assert.Empty(t, service.Expeditor())
	assert.Empty(t, service.From())

	err := service.Send(nil)
	assert.Error(t, err)

	err = service.Send(m)
	assert.Error(t, err)

	service.Config = &mail.Config{
		Host:      "smtp.world.com",
		Port:      "0",
		Username:  "secret",
		Password:  "secret",
		From:      "hello@world.com",
		Expeditor: "Whoami",
	}

	assert.NotEmpty(t, service.Expeditor())
	assert.NotEmpty(t, service.From())

	err = service.Send(m)
	assert.Error(t, err)

}
