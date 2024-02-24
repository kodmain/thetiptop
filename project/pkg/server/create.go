// Package server provides a simple HTTP/HTTPS server implementation for serving static and connect API.
package server

import (
	"fmt"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/swagger"
)

var server *Server

// Create return a instance os Server
// Pattern Singleton
func Create() *Server {
	if server == nil {
		config := fiber.Config{
			Prefork: true, // Multithreading
		}

		if os.Getppid() <= 1 {
			fmt.Println("WARNING: fiber in downgrade mode please use docker run --pid=host")
			config.Prefork = false // Disable to prevent bug in container
		}

		app := fiber.New(config)

		server = &Server{
			app:      app,
			api:      app.Group("api"),
			versions: make(map[string]fiber.Router),
		}

		server.app.Use(setGoToDoc)                           // register middleware setGoToDoc
		server.app.Use(setSecurityHeaders)                   // register middleware setSecurityHeaders
		server.app.Get("/docs/*", swagger.HandlerDefault)    // register middleware for documentation
		server.app.Group("/api", setRedirectOnEntryPointAPI) // entrypoint of the API but display we need to documentation
	}

	return server
}

// setGoToDoc is a middleware that redirect to /docs url path is like /
func setGoToDoc(c *fiber.Ctx) error {
	if c.Path() == "/index.html" || c.Path() == "/" {
		return c.Redirect("/docs", 301)
	}
	return c.Next()
}

// setRedirectOnEntryPointAPI is a middleware that redirect to /docs url path is like /api(?/)
func setRedirectOnEntryPointAPI(c *fiber.Ctx) error {
	if c.Path() == "/api" || c.Path() == "/api/" {
		return c.Redirect("/docs", 301)
	}
	return c.Next()
}

// setSecurityHeaders is a middleware that grants best practice around security
func setSecurityHeaders(c *fiber.Ctx) error {
	// Activer HSTS (HTTP Strict Transport Security)
	c.Set("Strict-Transport-Security", "max-age=63072000; includeSubDomains; preload")
	// Activer CSP (Content Security Policy)
	c.Set("Content-Security-Policy", "default-src 'unsafe-inline' 'self' fonts.gstatic.com fonts.googleapis.com;img-src data: 'self'")
	// Activer CORS (Cross-Origin Resource Sharing)
	c.Set("Access-Control-Allow-Origin", "*")
	c.Set("Access-Control-Allow-Methods", "GET,POST,HEAD,PUT,DELETE,PATCH")
	c.Set("Access-Control-Allow-Headers", "*")
	c.Set("Access-Control-Allow-Credentials", "true")

	return c.Next()
}
