// Package v1 grouping all API flagged by version number 1
package v1

import (
	"github.com/gofiber/fiber/v2"
	"github.com/kodmain/thetiptop/project/internal/api"
	"github.com/kodmain/thetiptop/project/internal/api/v1/fizzbuzz"
	"github.com/kodmain/thetiptop/project/internal/api/v1/status"
)

var (
	// Status is a variable that represents an instance of the `api.API` structure that defines a RESTful API for the "status" namespace with a single GET endpoint.
	// The endpoint is named "healthcheck" and it calls the `status.HealthCheck` function to handle the request.
	Status = &api.API{
		Version:   "v1",
		Namespace: "status",
		Get: map[string][]func(*fiber.Ctx) error{
			"healthcheck": {status.HealthCheck},
		},
	}

	// Fizzbuzz is a variable that represents an instance of the `api.API` structure that defines a RESTful API for the "fizzbuzz" namespace with two GET endpoints.
	// The first endpoint has the path `/:int1/:int2/:limit/:str1/:str2` and it calls three functions to handle the request: `fizzbuzz.FizzBuzzControls`, `fizzbuzz.FizzBuzzHits`, and `fizzbuzz.FizzBuzz`.
	// The second endpoint has the path "/stats" and it calls the `fizzbuzz.Stats` function to handle the request.
	Fizzbuzz = &api.API{
		Version:   "v1",
		Namespace: "fizzbuzz",
		Get: map[string][]func(*fiber.Ctx) error{
			"/:int1/:int2/:limit/:str1/:str2": {fizzbuzz.FizzBuzzControls, fizzbuzz.FizzBuzzHits, fizzbuzz.FizzBuzz},
			"/stats":                          {fizzbuzz.Stats},
		},
	}
)
