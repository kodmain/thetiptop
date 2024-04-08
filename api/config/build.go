// Package config provides functions/var/const for loading and accessing configuration settings for the application.
package config

var (
	BUILD_COMMIT  string               // The commit hash of the build, useful for tracking specific builds in version control.
	BUILD_VERSION string = "local"     // The version of the build, defaults to the value in DEFAULT_VERSION.
	APP_NAME      string = "TheTipTop" // The name of the application, defaults to the value in DEFAULT_APP_NAME.
	HOSTNAME      string = "localhost" // The hostname of the server, used for generating TLS certificates.
)
