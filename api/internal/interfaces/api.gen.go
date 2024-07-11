// Automatically generated by api/generator/api.gen.go, DO NOT EDIT manually
// Package api implements Register method for fiber
package interfaces

import (
	"encoding/json"

	"github.com/gofiber/fiber/v2"
	"github.com/kodmain/thetiptop/api/internal/docs"
	"github.com/kodmain/thetiptop/api/internal/interfaces/api/client"
	"github.com/kodmain/thetiptop/api/internal/interfaces/status"
	"github.com/swaggo/swag"
)

func init() {
	json.Unmarshal([]byte(doc), Mapping)
}

// API represents a collection of HTTP endpoints grouped by namespace and version.
var (
	Endpoints map[string]fiber.Handler = map[string]func(*fiber.Ctx) error{
		"client.PasswordRecover": client.PasswordRecover,
		"client.PasswordUpdate":  client.PasswordUpdate,
		"client.SignIn":          client.SignIn,
		"client.SignRenew":       client.SignRenew,
		"client.SignUp":          client.SignUp,
		"client.SignValidation":  client.SignValidation,
		"status.HealthCheck":     status.HealthCheck,
		"status.IP":              status.IP,
	}
	Mapping = &docs.Swagger{}
	doc, _  = swag.ReadDoc()
)
