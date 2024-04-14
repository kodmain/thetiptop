package config_test

import (
	"testing"

	"github.com/kodmain/thetiptop/api/config"
	"github.com/stretchr/testify/assert"
)

func TestLoad(t *testing.T) {
	assert.Nil(t, config.Get("mail"))
	err := config.Load("")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "path is required")

	err = config.Load("cnf.yml")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "open cnf.yml: no such file or directory")

	err = config.Load("../config.test.yml")
	assert.Nil(t, err)

	err = config.Load("s3://config.kodmain/config.yml")
	assert.Error(t, err) // This test will fail because no credentials are provided
}

func TestGet(t *testing.T) {
	config.Reset()
	assert.Nil(t, config.Get("mail"))
	config.Load("../config.test.yml")

	assert.Nil(t, config.Get("toto"))

	assert.NotNil(t, config.Get("mail"))
	assert.NotNil(t, config.Get("databases"))
	assert.NotNil(t, config.Get("jwt"))

	assert.NotNil(t, config.Get("mail.host"))
	assert.NotNil(t, config.Get("databases.default"))
	assert.NotNil(t, config.Get("jwt.secret"))
}
