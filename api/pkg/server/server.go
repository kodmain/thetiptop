// Package server provides a simple HTTP/HTTPS server implementation for serving static and connect API.
package server

import (
	"fmt"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/kodmain/thetiptop/api/config"
	"github.com/kodmain/thetiptop/api/internal/api"
	"github.com/kodmain/thetiptop/api/internal/lib"
)

var methods = []string{"GET", "POST", "PUT", "PATCH", "DELETE"}

// Server is a struct that represents a Fiber server instance with an underlying `fiber.App`, `fiber.Router` for the API endpoints, and a map of `fiber.Router` instances for different API versions.
type Server struct {
	// `app` is an instance of the `fiber.App` structure that represents the underlying Fiber server instance.
	app *fiber.App
	// `api` is an instance of the `fiber.Router` structure that represents the API endpoint router for the default API version.
	api fiber.Router
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
func (server *Server) Register(handlers map[string]fiber.Handler) {
	for url, pathItem := range api.Mapping.Paths {
		for _, method := range methods {
			if isMethodDefined(pathItem, method) {
				if handler, exists := handlers[pathItem.Get.OperationID]; exists {
					server.api.Add(method, url, handler)
					fmt.Println("Add handler", pathItem.Get.OperationID, method, url)
				} else {
					fmt.Println("Handler not found", pathItem.Get.OperationID)
				}
			}
		}
	}
}

func isMethodDefined(pathItem *api.PathItem, method string) bool {
	switch method {
	case "GET":
		return pathItem.Get != nil
	case "POST":
		return pathItem.Post != nil
	case "PUT":
		return pathItem.Put != nil
	case "PATCH":
		return pathItem.Patch != nil
	case "DELETE":
		return pathItem.Delete != nil
	}
	return false
}

func (server *Server) http() {
	lib.WithCriticalError(server.app.Listen(":80"))
}

func (server *Server) https() {
	lib.WithCriticalError(server.app.ListenTLS(":443", config.TLSPath+config.TLSCertFile, config.TLSPath+config.TLSKeyFile))
}
