// Package config provides functions/var/const for loading and accessing configuration settings for the application.
package config

import (
	"github.com/kodmain/thetiptop/api/internal/docs/generated"
)

// Initialize SwaggerInfo
func init() {
	generated.SwaggerInfo.Title = APP_NAME
	generated.SwaggerInfo.Description = APP_NAME + " API"
	generated.SwaggerInfo.Version = BUILD_VERSION
	generated.SwaggerInfo.Host = "localhost"
	generated.SwaggerInfo.BasePath = ""
	generated.SwaggerInfo.Schemes = []string{}
}
