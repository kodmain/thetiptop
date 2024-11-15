package env

import (
	"os"
	"strings"
)

var (
	BUILD_COMMIT   string                 // The commit hash of the build, useful for tracking specific builds in version control.
	BUILD_VERSION  string                 // The version of the build, defaults to the value in DEFAULT_VERSION.
	APP_NAME       string = "TheTipTop"   // The name of the application, defaults to the value in DEFAULT_APP_NAME.
	APP_LABEL_NAME string = "The Tip Top" // The name of the application, defaults to the value in DEFAULT_APP_NAME.
	HOSTNAME       string = "localhost"   // The hostname of the server, used for generating TLS certificates.

	//DEFAULT_DB_NAME     string = "default"
	DEFAULT_CONFIG_URI  string = "s3://config.kodmain/config.yml"
	DEFAULT_AWS_PROFILE string = ""
	DEFAULT_PORT_HTTP   int    = 80
	DEFAULT_PORT_HTTPS  int    = 443

	CONFIG_URI  *string = &DEFAULT_CONFIG_URI
	AWS_PROFILE *string = &DEFAULT_AWS_PROFILE
	PORT_HTTP   *int    = &DEFAULT_PORT_HTTP
	PORT_HTTPS  *int    = &DEFAULT_PORT_HTTPS
)

func IsTest() bool {
	return strings.HasSuffix(os.Args[0], ".test")
}
