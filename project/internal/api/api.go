// Package api implements Register method for fiber
package api

import (
	"reflect"

	"github.com/gofiber/fiber/v2"
)

var (
	methods = []string{"Get", "Post", "Put", "Patch", "Delete"}
)

// API represents a collection of HTTP endpoints grouped by namespace and version.
type API struct {
	// Namespace specifies the namespace of the API (e.g. "user" or "product").
	Namespace string
	// Version specifies the version of the API (e.g. "v1" or "v2").
	Version string
	// Get is a map of HTTP GET endpoints, where the keys are the endpoint paths
	// (e.g. "/users") and the values are arrays of handler functions to be executed
	// for each request to that endpoint.
	Get map[string][]fiber.Handler
	// Post is a map of HTTP POST endpoints, with the same structure as Get.
	Post map[string][]fiber.Handler
	// Put is a map of HTTP PUT endpoints, with the same structure as Get.
	Put map[string][]fiber.Handler
	// Patch is a map of HTTP PATCH endpoints, with the same structure as Get.
	Patch map[string][]fiber.Handler
	// Delete is a map of HTTP DELETE endpoints, with the same structure as Get.
	Delete map[string][]fiber.Handler
}

// Register registers the API's endpoints with the given router.
// It uses reflection to iterate over the API's HTTP endpoint maps (Get, Post, Put, Patch, Delete),
// extracts the endpoint paths and handler functions, and adds them to the router using the Add method.
// The router is typically an instance of the fiber.Router type.
func (api *API) Register(router fiber.Router) {
	apiValue := reflect.ValueOf(*api)
	for _, method := range methods {
		methodValue := apiValue.FieldByName(method)
		for _, key := range methodValue.MapKeys() {
			route := key.String()
			handlersValue := methodValue.MapIndex(key)
			handlers := handlersValue.Interface().([]fiber.Handler)
			router.Add(method, route, handlers...)
		}
	}
}
