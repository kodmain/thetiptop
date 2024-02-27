package config

var (
	BUILD_COMMIT  string // The commit hash of the build, useful for tracking specific builds in version control.
	BUILD_VERSION string // The version of the build, defaults to the value in DEFAULT_VERSION.
)
