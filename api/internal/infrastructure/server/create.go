// Package server provides a simple HTTP/HTTPS server implementation for serving static and connect API.
package server

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/swagger"
	"github.com/kodmain/thetiptop/api/config"
	"github.com/kodmain/thetiptop/api/internal/docs/generated"
	"github.com/kodmain/thetiptop/api/internal/infrastructure/server/certs"
)

var servers map[string]*Server = make(map[string]*Server)

func getConfig(cfgs ...fiber.Config) fiber.Config {
	if len(cfgs) > 0 {
		return cfgs[0]
	}

	cfg := fiber.Config{
		AppName:               config.APP_NAME,
		DisableStartupMessage: true,  // Disable Prefork to prevent bug in container and because SO_REUSEPORT can give false metrics in prometheus, maybe in the future we can use REDIS to store metrics
		Prefork:               false, // Disable multithreading
	}

	return cfg
}

// Create return a instance os Server
// Pattern Singleton
func Create(cfgs ...fiber.Config) *Server {
	cfg := getConfig(cfgs...)
	if server, exists := servers[cfg.AppName]; exists {
		return server
	}

	server := &Server{
		app:   fiber.New(cfg),
		certs: certs.TLSConfigFor(config.HOSTNAME),
	}

	server.app.Use(setGoToDoc)         // register middleware setGoToDoc
	server.app.Use(setSecurityHeaders) // register middleware setSecurityHeaders
	server.app.Get("/docs/*", swagger.New(swagger.Config{
		Title:                    config.APP_NAME,
		Layout:                   "BaseLayout",
		DocExpansion:             "list",
		DefaultModelsExpandDepth: 2,
	})) // register middleware for documentation

	servers[cfg.AppName] = server

	return server
}

// setGoToDoc is a middleware that redirect to /docs url path is like /
func setGoToDoc(c *fiber.Ctx) error {
	if c.Path() == "/index.html" || c.Path() == "/" {
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

	generated.SwaggerInfo.Host = c.Hostname()

	return c.Next()
}
