// Package server provides a simple HTTP/HTTPS server implementation for serving static and connect API.
package server

import (
	"crypto/tls"
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/kodmain/thetiptop/api/env"
	"github.com/kodmain/thetiptop/api/internal/application"
	"github.com/kodmain/thetiptop/api/internal/docs"
	"github.com/kodmain/thetiptop/api/internal/infrastructure/observability/logger"
	"github.com/kodmain/thetiptop/api/internal/interfaces"
)

var methods = []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS", "HEAD"}

const (
	v   = "%v"
	vv  = "%v %v"
	vvv = "%v %v %v"
)

// Server is a struct that represents a Fiber server instance with an underlying `fiber.App`, `fiber.Router` for the API endpoints, and a map of `fiber.Router` instances for different API versions.
type Server struct {
	// `app` is an instance of the `fiber.App` structure that represents the underlying Fiber server instance.
	app *fiber.App
	// `certs` is a pointer to a `tls.Config` structure that holds the TLS configuration for the server.
	certs *tls.Config
}

// Start is a method of the `Server` struct that starts the server and listens for incoming HTTP and/or HTTPS requests, depending on the `config` settings.
func (server *Server) Start() error {
	// Launch a goroutine to start the HTTP server.
	go server.http()
	// If HTTPS is enabled in the `config`, launch a goroutine to start the HTTPS server as well.
	go server.https()

	return nil
}

func (server *Server) Stop() error {
	logger.Info("Shutting down server")
	return server.app.Shutdown()
}

// API is a method of the `Server` struct that registers the provided API with the server.
// It creates a new version of the API router, adds the provided API to the router's namespace, and registers the new router with the server's main `app` instance.
func (server *Server) Register(handlers map[string]fiber.Handler) {
	for url, pathItem := range interfaces.Mapping.Paths {
		// Remplacer {property} par :property dans l'URL
		url = replaceCurlyBracesWithColons(url)

		for _, method := range methods {
			if operationID, exist := getOperationID(pathItem, method); exist {
				handlersToRegister := getHandlers(operationID, handlers)
				if len(handlersToRegister) > 0 {
					server.app.Add(method, url, wrapLastHandler(handlersToRegister)...)
					logger.Infof("Register %v %v with %v", method, url, operationID)
				} else {
					logger.Warnf("Handler not found %v", operationID)
				}
			}
		}
	}
}

// replaceCurlyBracesWithColons remplace toutes les occurences {property} par :property
func replaceCurlyBracesWithColons(url string) string {
	// Utilisation d'une expression régulière pour remplacer {property} par :property
	re := regexp.MustCompile(`\{([a-zA-Z0-9_]+)\}`)
	return re.ReplaceAllString(url, ":$1")
}

func wrapLastHandler(handlers []fiber.Handler) []fiber.Handler {
	if len(handlers) == 0 {
		return handlers
	}

	last := handlers[len(handlers)-1]
	wrappedLast := func(c *fiber.Ctx) error {
		err := last(c)
		status := c.Response().StatusCode()
		if status >= 400 {
			logger.Error(fmt.Errorf(vvv, c.Method(), status, c.Path()))
		} else {
			logger.Messagef(vvv, c.Method(), c.Response().StatusCode(), c.Path())
		}

		return err
	}

	return append(handlers[:len(handlers)-1], wrappedLast)
}

func getHandlers(operationID string, handlers map[string]fiber.Handler) []fiber.Handler {
	var handlersToRegister []fiber.Handler
	var operationIDs = strings.Split(operationID, "=>")
	for key, handler := range operationIDs {
		handler = strings.TrimSpace(handler)
		if h, exist := handlers[handler]; exist {
			handlersToRegister = append(handlersToRegister, h)
		} else {
			logger.Warnf("Handler not found %v %v", key, handler)
		}
	}
	return handlersToRegister
}

func getOperationID(pathItem *docs.PathItem, method string) (string, bool) {
	switch method {
	case "GET":
		if pathItem.Get != nil {
			return pathItem.Get.OperationID, true
		}
	case "POST":
		if pathItem.Post != nil {
			return pathItem.Post.OperationID, true
		}
	case "PUT":
		if pathItem.Put != nil {
			return pathItem.Put.OperationID, true
		}
	case "DELETE":
		if pathItem.Delete != nil {
			return pathItem.Delete.OperationID, true
		}
	case "PATCH":
		if pathItem.Patch != nil {
			return pathItem.Patch.OperationID, true
		}
	case "OPTIONS":
		if pathItem.Options != nil {
			return pathItem.Options.OperationID, true
		}
	case "HEAD":
		if pathItem.Head != nil {
			return pathItem.Head.OperationID, true
		}
	}
	return "", false
}

func (server *Server) http() {
	logger.Info("Server http started on port " + strconv.Itoa(*env.PORT_HTTP))
	application.PANIC <- server.app.Listen(":" + strconv.Itoa(*env.PORT_HTTP))
}

func (server *Server) https() {
	logger.Info("Server https started on port " + strconv.Itoa(*env.PORT_HTTPS))
	application.PANIC <- server.app.ListenTLSWithCertificate(":"+strconv.Itoa(*env.PORT_HTTPS), server.certs.Certificates[0])
}
