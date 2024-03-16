package status_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/kodmain/thetiptop/api/internal/api/status"
	"github.com/stretchr/testify/assert"
)

func TestIP(t *testing.T) {
	app := fiber.New()
	assert.NotNil(t, app)

	app.Get("/status/ip", status.IP)

	req := httptest.NewRequest(http.MethodGet, "/status/ip", nil)
	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}
