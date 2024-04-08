package server

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
)

func TestCreateServer(t *testing.T) {
	srvA := Create()
	assert.NotNil(t, srvA)
	srvB := Create()
	assert.NotNil(t, srvB)

	assert.Equal(t, srvA, srvB)
}

func TestSetGoToDoc(t *testing.T) {
	app := fiber.New()
	app.Use(setGoToDoc)
	app.Get("/hello", func(c *fiber.Ctx) error {
		return c.SendString("World!")
	})

	req := httptest.NewRequest(http.MethodGet, "/index.html", nil)
	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusMovedPermanently, resp.StatusCode)
	assert.Equal(t, "/docs", resp.Header.Get("Location"))

	req = httptest.NewRequest(http.MethodGet, "/hello", nil)
	resp, err = app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	data, err := io.ReadAll(resp.Body)
	assert.NoError(t, err)
	defer resp.Body.Close()
	assert.Equal(t, "World!", string(data))

}

func TestSetSecurityHeaders(t *testing.T) {
	app := fiber.New()
	app.Use(setSecurityHeaders)

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusNotFound, resp.StatusCode)
	assert.Equal(t, "max-age=63072000; includeSubDomains; preload", resp.Header.Get("Strict-Transport-Security"))
	assert.Equal(t, "default-src 'unsafe-inline' 'self' fonts.gstatic.com fonts.googleapis.com;img-src data: 'self'", resp.Header.Get("Content-Security-Policy"))
	assert.Equal(t, "*", resp.Header.Get("Access-Control-Allow-Origin"))
	assert.Equal(t, "GET,POST,HEAD,PUT,DELETE,PATCH", resp.Header.Get("Access-Control-Allow-Methods"))
	assert.Equal(t, "*", resp.Header.Get("Access-Control-Allow-Headers"))
	assert.Equal(t, "true", resp.Header.Get("Access-Control-Allow-Credentials"))
}
