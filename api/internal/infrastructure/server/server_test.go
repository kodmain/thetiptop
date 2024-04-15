package server

import (
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/kodmain/thetiptop/api/config"
	"github.com/kodmain/thetiptop/api/internal/interfaces/status"
	"github.com/stretchr/testify/assert"
)

func TestRegister(t *testing.T) {
	srv := Create(fiber.Config{
		AppName: config.APP_NAME,
		Prefork: false, // Multithreading
	})

	srv.Register(map[string]fiber.Handler{
		"status.HealthCheck": status.HealthCheck,
	})

	assert.NotNil(t, srv)
	assert.Nil(t, srv.Start())
	assert.Nil(t, srv.Stop())
}
