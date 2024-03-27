package config

import "github.com/kodmain/thetiptop/api/internal/application/observability/logger/levels"

var (
	// DEFAULT_LOG_LEVEL defines the default log level for the application.
	// This level is used to control the verbosity of the application's logs.
	// It can be set to any of the levels defined in the levels package.
	// By default, it is set to levels.INFO.
	// To change the default log level, set the DEFAULT_LOG_LEVEL constant in the config package.
	// For example:
	// config.DEFAULT_LOG_LEVEL = levels.DEBUG
	// config.DEFAULT_LOG_LEVEL = levels.ERROR
	// config.DEFAULT_LOG_LEVEL = levels.FATAL
	// config.DEFAULT_LOG_LEVEL = levels.INFO
	// config.DEFAULT_LOG_LEVEL = levels.TRACE
	DEFAULT_LOG_LEVEL levels.TYPE = levels.DEFAULT
)
