// Package fizzbuzz implementation all handler for FizzBuzz API
package fizzbuzz

import (
	"sync"

	"github.com/gofiber/fiber/v2"
	"github.com/kodmain/thetiptop/project/internal/model/fizzbuzz"
)

var statistics = make(map[fizzbuzz.Request]uint)
var mutex = &sync.Mutex{}

// Stats is an HTTP handler function that returns statistics about the most frequently requested
// FizzBuzz endpoint. It retrieves the statistics from a shared global map that is updated by the
// "FizzBuzzHits" function. The statistics include the FizzBuzz request parameters that were used
// most frequently, as well as the number of times that those parameters were used.
//
// Parameters:
// - c: a pointer to the fiber.Ctx object representing the HTTP request context.
//
// Returns:
// - an error value, which is typically nil, since there is no meaningful error condition for this endpoint.
// @Summary		Return FizzBuzz statistics.
// @Description	Return the parameters corresponding to the most used request, as well as the number of hits for this request.
// @Tags		FizzBuzz
// @Accept		*/*
// @Produce		json
// @Help		Name    ?        type    required   description
// @Success		200		{object}  fizzbuzz.Stats
// @Router		/api/v1/fizzbuzz/stats [get]
func Stats(c *fiber.Ctx) error {
	mutex.Lock()
	defer mutex.Unlock()

	var mostUsedRequest fizzbuzz.Request
	var maxHits uint

	for req, hits := range statistics {
		if hits > maxHits {
			maxHits = hits
			mostUsedRequest = req
		}
	}

	return c.JSON(fizzbuzz.Stats{
		Request: mostUsedRequest,
		Hits:    maxHits,
	})
}

// FizzBuzzHits is an HTTP middleware function that updates a global map of FizzBuzz request statistics
// with information about the current request. It retrieves the FizzBuzz request parameters from the
// user value stored in the request context, increments the corresponding hit counter in the statistics
// map, and then passes the request on to the next middleware or endpoint handler.
//
// Parameters:
// - c: a pointer to the fiber.Ctx object representing the HTTP request context.
//
// Returns:
// - an error value, which is typically nil, since there is no meaningful error condition for this endpoint.
func FizzBuzzHits(c *fiber.Ctx) error {
	mutex.Lock()
	defer mutex.Unlock()

	var fbr fizzbuzz.Request = c.Context().UserValue("fizzbuzz.Request").(fizzbuzz.Request)
	statistics[fbr]++

	return c.Next()
}
