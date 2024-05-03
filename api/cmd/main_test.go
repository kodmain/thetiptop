package main

import (
	"bytes"
	"os"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/kodmain/thetiptop/api/env"
	"github.com/kodmain/thetiptop/api/internal/application"
	"github.com/stretchr/testify/assert"
)

// TestHelperPreRunE tests the PreRunE function for configuration loading
func TestHelperPreRunE(t *testing.T) {
	env.CONFIG_URI = aws.String("../config.test.yml")
	cmd := Helper
	b := bytes.NewBufferString("")
	cmd.SetOut(b)
	cmd.SetErr(b)
	err := cmd.PreRunE(cmd, nil)
	assert.Nil(t, err)
}

func TestHelperRunE(t *testing.T) {
	env.PORT_HTTP = aws.Int(8080)
	env.PORT_HTTPS = aws.Int(8443)

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

func TestVersionCmd(t *testing.T) {
	cmd := versionCmd
	b := bytes.NewBufferString("")
	cmd.SetOut(b)
	cmd.SetErr(b)
	cmd.Run(cmd, nil)
}
