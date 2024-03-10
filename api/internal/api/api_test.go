package api

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
)

func TestRegister(t *testing.T) {
	// Create a new instance of API
	const Users = "/users"
	const Api = "/api" + Users
	apiTest := &API{
		Namespace: "test",
		Version:   "v1",
		Get: map[string][]fiber.Handler{
			Users: {func(c *fiber.Ctx) error {
				return c.SendString("GET /users")
			}},
		},
		Post: map[string][]fiber.Handler{
			Users: {func(c *fiber.Ctx) error {
				return c.SendString("POST /users")
			}},
		},
		Put: map[string][]fiber.Handler{
			Users + "/:id": {func(c *fiber.Ctx) error {
				return c.SendString("PUT /users/:id")
			}},
		},
		Patch: map[string][]fiber.Handler{
			Users + "/:id": {func(c *fiber.Ctx) error {
				return c.SendString("PATCH /users/:id")
			}},
		},
		Delete: map[string][]fiber.Handler{
			Users + "/:id": {func(c *fiber.Ctx) error {
				return c.SendString("DELETE /users/:id")
			}},
		},
	}

	// Create a new router and register the API
	router := fiber.New()
	apiTest.Register(router.Group("api"))

	// Test each endpoint using HTTP requests
	testCases := []struct {
		method string
		path   string
		body   string
		want   string
	}{
		{"GET", Api, "", "GET /users"},
		{"POST", Api, "", "POST /users"},
		{"PUT", Api + "/123", "", "PUT /users/:id"},
		{"PATCH", Api + "/123", "", "PATCH /users/:id"},
		{"DELETE", Api + "/123", "", "DELETE /users/:id"},
	}

	for _, tc := range testCases {
		req := httptest.NewRequest(tc.method, tc.path, nil)
		resp, err := router.Test(req, -1)
		if err != nil {
			t.Fatalf("TestAPI_Register: %v", err)
		}
		if resp.StatusCode != http.StatusOK {
			t.Errorf("TestAPI_Register: unexpected status code %d", resp.StatusCode)
		}
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			t.Fatalf("TestAPI_Register: %v", err)
		}
		if string(body) != tc.want {
			t.Errorf("TestAPI_Register: unexpected response body: %q", string(body))
		}
	}
}
