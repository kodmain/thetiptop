package mail_test

import (
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/kodmain/thetiptop/api/config"
	"github.com/kodmain/thetiptop/api/internal/infrastructure/providers/mail"
	"github.com/stretchr/testify/assert"
)

const (
	GOOD_EMAIL = "user1@example.com"
	GOOD_PASS  = "ValidP@ssw0rd1"
)

func TestMail(t *testing.T) {
	config.Load(aws.String("../../../../config.test.yml"))

	tpl := mail.NewTemplate("signup")
	text, html, err := tpl.Inject(mail.Data{
		"AppName": "ThéTipTop",
		"Url":     "https://thetiptop.com",
	})

	assert.NoError(t, err)

	m := &mail.Mail{
		To:      []string{GOOD_EMAIL},
		Subject: "Welcome to The Tip Top",
		Text:    text,
		Html:    html,
	}

	assert.True(t, m.IsValid())

	assert.NotNil(t, m)
	msg, to, err := m.Prepare()
	assert.NotNil(t, msg)
	assert.NotNil(t, to)
	assert.NoError(t, err)
}
