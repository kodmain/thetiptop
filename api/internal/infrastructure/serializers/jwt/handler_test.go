package jwt_test

import (
	"io"
	"net/http"
	"net/url"
	"strings"
	"testing"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/kodmain/thetiptop/api/internal/infrastructure/serializers/buffer"
	"github.com/kodmain/thetiptop/api/internal/infrastructure/serializers/jwt"
	"github.com/stretchr/testify/assert"
)

var fbr *fiber.App = fiber.New()

func start() error {
	jwt.New(&jwt.JWT{
		Expire:   5,
		Refresh:  10,
		Duration: time.Second,
	})

	// Add your authentication middleware to routes that require authentication
	fbr.Use(jwt.Parser)

	// Define your route handlers here
	fbr.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello, World!")
	})

	// Define your route handlers here
	fbr.Get("/restricted", jwt.Auth, func(c *fiber.Ctx) error {
		return c.SendString("Hello, Restricted!")
	})

	c := make(chan error, 1)

	time.AfterFunc(1*time.Second, func() {
		c <- nil
	})

	go func() {
		c <- fbr.Listen(":3000")
	}()

	// Start the server on port 3000
	return <-c
}

func stop() error {
	return fbr.Shutdown()
}

func request(method, uri string, token string, values ...map[string][]string) ([]byte, int, error) {
	// Create a form with email and password fields
	form := url.Values{}
	if len(values) > 0 {
		for key, value := range values[0] {
			form.Set(key, value[0])
		}
	}

	// Create a new HTTP request to call /sign/up
	req, err := http.NewRequest(method, uri, nil)
	if err != nil {
		return nil, 0, err
	}

	// Set the form as the request body
	if len(values) > 0 {
		req.Body = io.NopCloser(strings.NewReader(form.Encode()))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	}

	if token != "" {
		req.Header.Set("Authorization", token)
	}

	// Perform the request
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, 0, err
	}

	defer resp.Body.Close()

	content, err := buffer.Read(resp.Body)

	return content.Bytes(), resp.StatusCode, err
}

func TestParser(t *testing.T) {
	assert.NoError(t, start())

	const (
		restricted = "http://localhost:3000/restricted"
		bearer     = "Bearer "
	)

	t.Run("TestHomePage", func(t *testing.T) {
		content, status, err := request("GET", "http://localhost:3000", "", nil)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, status)
		assert.Equal(t, "Hello, World!", string(content))
	})

	t.Run("TestRestrictedNoToken", func(t *testing.T) {
		content, status, err := request("GET", restricted, "", nil)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusUnauthorized, status)
		assert.Equal(t, "{\"code\":401,\"message\":\"auth.no_token\"}", string(content))
	})

	t.Run("TestTokenCreation", func(t *testing.T) {
		token, refresh, err := jwt.FromID("hello", nil)
		assert.NoError(t, err)
		assert.NotEmpty(t, token)
		assert.NotEmpty(t, refresh)
	})

	t.Run("TestRestrictedBadFormat", func(t *testing.T) {
		token, _, _ := jwt.FromID("hello", nil)
		content, status, err := request("GET", restricted, token, nil)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, status)
		assert.Equal(t, "{\"code\":400,\"message\":\"auth.bad_format\"}", string(content))
	})

	t.Run("TestRestrictedInvalidToken", func(t *testing.T) {
		token, _, _ := jwt.FromID("hello", nil)
		content, status, err := request("GET", restricted, bearer+"Oki"+token, nil)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusUnauthorized, status)
		assert.Equal(t, "{\"code\":401,\"message\":\"auth.failed\"}", string(content))
	})

	t.Run("TestRestrictedValidToken", func(t *testing.T) {
		token, _, _ := jwt.FromID("hello", nil)
		content, status, err := request("GET", restricted, bearer+token, nil)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, status)
		assert.Equal(t, "Hello, Restricted!", string(content))
	})

	t.Run("TestRestrictedExpiredToken", func(t *testing.T) {
		token, _, _ := jwt.FromID("hello", nil)
		time.Sleep(5 * time.Second)
		content, status, err := request("GET", restricted, bearer+token, nil)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusUnauthorized, status)
		assert.Equal(t, "{\"code\":401,\"message\":\"auth.failed\"}", string(content))
	})

	t.Run("TestNewSecretInvalidatesOldToken", func(t *testing.T) {
		realToken, refreshToken, err := jwt.FromID("hello", nil)
		assert.NoError(t, err)
		assert.NotEmpty(t, realToken)
		assert.NotEmpty(t, refreshToken)

		jwt.New(&jwt.JWT{
			Secret: "secret",
		})

		content, status, errhttp := request("GET", restricted, bearer+realToken, nil)
		assert.NoError(t, errhttp)
		assert.Equal(t, http.StatusUnauthorized, status)
		assert.Equal(t, "{\"code\":401,\"message\":\"auth.failed\"}", string(content))
	})

	assert.NoError(t, stop())
}
