// Package config provides functions/var/const for loading and accessing configuration settings for the application.
package config

import (
	"github.com/kodmain/thetiptop/api/internal/docs"
)

// Initialize SwaggerInfo
func init() {
	docs.SwaggerInfo.Title = APP_NAME
	docs.SwaggerInfo.Description = "TheTipTop API"
	docs.SwaggerInfo.Version = BUILD_VERSION
	docs.SwaggerInfo.Host = "localhost"
	docs.SwaggerInfo.BasePath = ""
	docs.SwaggerInfo.Schemes = []string{}
}
