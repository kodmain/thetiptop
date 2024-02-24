package server

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/kodmain/thetiptop/project/internal/api"
	"github.com/stretchr/testify/assert"
)

func TestServerAPI(t *testing.T) {
	server := &Server{
		app:      fiber.New(),
		api:      fiber.New(),
		versions: make(map[string]fiber.Router),
	}

	api := &api.API{
		Version:   "v1",
		Namespace: "/test",
	}

	server.API(api)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/test", nil)
	resp, err := server.app.Test(req)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusNotFound, resp.StatusCode)
}
