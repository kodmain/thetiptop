package fizzbuzz_test

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/kodmain/thetiptop/api/internal/api/fizzbuzz"
	"github.com/stretchr/testify/assert"
)

const (
	url                  = "/fizzbuzz/:int1/:int2/:limit/:str1/:str2"
	shouldNotReturnError = "Request should not return an error"
	shouldBe200          = "Response status code should be 200"
)

func TestFizzBuzz(t *testing.T) {
	tests := []struct {
		name         string
		url          string
		expectedCode int
		expectedBody string
	}{
		{
			name:         "Valid parameters",
			url:          "/fizzbuzz/3/5/15/Fizz/Buzz",
			expectedCode: fiber.StatusOK,
			expectedBody: `["1","2","Fizz","4","Buzz","Fizz","7","8","Fizz","Buzz","11","Fizz","13","14","FizzBuzz"]`,
		},
		{
			name:         "Invalid int1",
			url:          "/fizzbuzz/invalid/invalid/invalid/Fizz/Buzz",
			expectedCode: fiber.StatusBadRequest,
			expectedBody: `{"error":{"different":"int1 and int2 must be different","int1":"invalid int1","int2":"invalid int2","limit":"invalid limit","limit_range":"limit must be between 1 and 100","no_zero":"int1, int2 and limit must be greater than 0"}}`,
		},
	}

	app := fiber.New()
	app.Get(url, fizzbuzz.FizzBuzz)

	for _, tt := range tests {
		req := httptest.NewRequest(http.MethodGet, tt.url, nil)
		resp, err := app.Test(req)

		assert.NoError(t, err, shouldNotReturnError)
		assert.Equal(t, tt.expectedCode, resp.StatusCode, "bad status code")

		data, err := io.ReadAll(resp.Body)
		assert.NoError(t, err)
		defer resp.Body.Close()
		assert.Equal(t, tt.expectedBody, string(data))
	}
}
