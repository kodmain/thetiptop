// Package server provides a simple HTTP/HTTPS server implementation for serving static and connect API.
package server

import (
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/kodmain/thetiptop/api/config"
	"github.com/kodmain/thetiptop/api/internal/api"
	"github.com/kodmain/thetiptop/api/internal/lib"
)

// Server is a struct that represents a Fiber server instance with an underlying `fiber.App`, `fiber.Router` for the API endpoints, and a map of `fiber.Router` instances for different API versions.
type Server struct {
	// `app` is an instance of the `fiber.App` structure that represents the underlying Fiber server instance.
	app *fiber.App
	// `api` is an instance of the `fiber.Router` structure that represents the API endpoint router for the default API version.
	api fiber.Router
	// `versions` is a map of `fiber.Router` instances that represent the API endpoint routers for different API versions.
	versions map[string]fiber.Router
}

// Start is a method of the `Server` struct that starts the server and listens for incoming HTTP and/or HTTPS requests, depending on the `config` settings.
func (server *Server) Start() {
	// Launch a goroutine to start the HTTP server.
	go server.http()
	// If HTTPS is enabled in the `config`, launch a goroutine to start the HTTPS server as well.
	if config.EnableHTTPS {
		go server.https()
	}
	// Wait for the server to exit, and exit the process with the appropriate status code.
	os.Exit(lib.WaitStatus())
}

// API is a method of the `Server` struct that registers the provided API with the server.
// It creates a new version of the API router, adds the provided API to the router's namespace, and registers the new router with the server's main `app` instance.
func (server *Server) API(api *api.API) {
	api.Register(server.version(api).Group(api.Namespace))
}

func (server *Server) version(api *api.API) fiber.Router {
	if router, exist := server.versions[api.Version]; exist {
		return router
	}

	server.versions[api.Version] = server.api.Group(api.Version)

	return server.versions[api.Version]
}

func (server *Server) http() {
	lib.WithCriticalError(server.app.Listen(":80"))
}

func (server *Server) https() {
	lib.WithCriticalError(server.app.ListenTLS(":443", config.TLSPath+config.TLSCertFile, config.TLSPath+config.TLSKeyFile))
}
