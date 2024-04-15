package mail_test

import (
	"testing"

	"github.com/kodmain/thetiptop/api/internal/infrastructure/providers/mail"
	"github.com/stretchr/testify/assert"
)

func TestTemplate(t *testing.T) {
	template := mail.NewTemplate("signup")
	assert.NotNil(t, template)

	text, html, err := template.Inject(mail.Data{
		"AppName": "ThéTipTop",
		"Url":     "https://thetiptop.com",
	})

	assert.NoError(t, err)
	assert.NotEmpty(t, text)
	assert.NotEmpty(t, html)

	template = mail.NewTemplate("oki")
	assert.Nil(t, template)

	err = mail.New(nil)
	assert.NotNil(t, err)

	err = mail.Send(&mail.Mail{
		To:      []string{"hello@world.com"},
		Subject: "Subject",
		Text:    text,
		Html:    html,
	})

	assert.Error(t, err)

	err = mail.New(&mail.Service{
		Disable: true,
	})

	assert.NotNil(t, err)

	err = mail.Send(&mail.Mail{})
	assert.Error(t, err)

	mail.Send(&mail.Mail{
		To:      []string{"hello@world.com"},
		Subject: "Subject",
		Text:    text,
		Html:    html,
	})
}
