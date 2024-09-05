package template_test

import (
	"testing"

	"github.com/kodmain/thetiptop/api/internal/infrastructure/providers/mail/template"
	"github.com/stretchr/testify/assert"
)

func TestData(t *testing.T) {
	tpl := template.NewTemplate("token")
	assert.NotNil(t, tpl)

	text, html, err := tpl.Inject(template.Data{
		"AppName": "ThéTipTop",
		"Url":     "https://thetiptop.com",
	})

	assert.NoError(t, err)
	assert.NotEmpty(t, text)
	assert.NotEmpty(t, html)

	tpl = template.NewTemplate("oki")
	assert.Nil(t, tpl)
}
