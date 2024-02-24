// Package config provides functions/var/const for loading and accessing configuration settings for the application.
package config

import "github.com/kodmain/thetiptop/project/internal/docs"

// Initialize SwaggerInfo
func init() {
	docs.SwaggerInfo.Title = "FizzBuzz"
	docs.SwaggerInfo.Description = "FizzBuzz API"
	docs.SwaggerInfo.Version = "1.0"
	docs.SwaggerInfo.Host = "localhost"
	docs.SwaggerInfo.BasePath = "/"
	docs.SwaggerInfo.Schemes = []string{"http", "https"}
}
