package metrics_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/kodmain/thetiptop/api/internal/api/metrics"
	"github.com/stretchr/testify/assert"
)

const (
	url                  = "/metrics/statistics"
	shouldNotReturnError = "Request should not return an error"
	shouldBe200          = "Response status code should be 200"
)

func TestStatistics(t *testing.T) {
	app := fiber.New()
	app.Get(url, metrics.Statistics)

	req := httptest.NewRequest(http.MethodGet, url, nil)
	resp, err := app.Test(req)
	assert.NoError(t, err, shouldNotReturnError)
	assert.Equal(t, http.StatusOK, resp.StatusCode, shouldBe200)

	app.Use(metrics.Counter)
	req = httptest.NewRequest(http.MethodGet, "/test", nil)
	resp, err = app.Test(req)
	assert.NoError(t, err, shouldNotReturnError)
	assert.Equal(t, http.StatusNotFound, resp.StatusCode, shouldBe200)

	req = httptest.NewRequest(http.MethodGet, url, nil)
	resp, err = app.Test(req)
	assert.NoError(t, err, shouldNotReturnError)
	assert.Equal(t, http.StatusOK, resp.StatusCode, shouldBe200)
}
