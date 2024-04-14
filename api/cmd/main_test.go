package main

import (
	"bytes"
	"os"
	"testing"
	"time"

	"github.com/kodmain/thetiptop/api/config"
	"github.com/kodmain/thetiptop/api/internal/application"
	"github.com/stretchr/testify/assert"
)

// TestHelperPreRunE tests the PreRunE function for configuration loading
func TestHelperPreRunE(t *testing.T) {
	config.DEFAULT_CONFIG = "../config.test.yml"
	cmd := Helper
	b := bytes.NewBufferString("")
	cmd.SetOut(b)
	cmd.SetErr(b)
	err := cmd.PreRunE(cmd, nil)
	assert.Nil(t, err)
}

func TestHelperRunE(t *testing.T) {
	config.PORT_HTTP = ":8080"
	config.PORT_HTTPS = ":8443"

	cmd := Helper
	b := bytes.NewBufferString("")
	cmd.SetOut(b)
	cmd.SetErr(b)
	err := cmd.RunE(cmd, nil)
	assert.Nil(t, err)
}

func TestHelperPostRunE(t *testing.T) {
	cmd := Helper
	b := bytes.NewBufferString("")
	cmd.SetOut(b)
	cmd.SetErr(b)
	time.AfterFunc(1*time.Second, func() {
		application.SIGS <- os.Interrupt
	})

	err := cmd.PostRunE(cmd, nil)
	assert.Nil(t, err)
}
