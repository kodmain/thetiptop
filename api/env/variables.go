package env

import (
	"flag"
)

var (
	BUILD_COMMIT  string               // The commit hash of the build, useful for tracking specific builds in version control.
	BUILD_VERSION string = "dev"       // The version of the build, defaults to the value in DEFAULT_VERSION.
	APP_NAME      string = "TheTipTop" // The name of the application, defaults to the value in DEFAULT_APP_NAME.
	HOSTNAME      string = "localhost" // The hostname of the server, used for generating TLS certificates.

	//DEFAULT_DB_NAME     string = "default"
	DEFAULT_CONFIG_URI  string = "s3://config.kodmain/config.yml"
	DEFAULT_AWS_PROFILE string = "kodmain"
	DEFAULT_PORT_HTTP   string = ":80"
	DEFAULT_PORT_HTTPS  string = ":443"

	//DB_NAME     = flag.String("dbname", DEFAULT_DB_NAME, "Nom de la base de données par défaut")
	CONFIG_URI  = flag.String("config", DEFAULT_CONFIG_URI, "URI de la configuration")
	AWS_PROFILE = flag.String("aws-profile", DEFAULT_AWS_PROFILE, "Profil AWS")
	PORT_HTTP   = flag.String("port-http", DEFAULT_PORT_HTTP, "Port HTTP")
	PORT_HTTPS  = flag.String("port-https", DEFAULT_PORT_HTTPS, "Port HTTPS")
)

func init() {
	flag.Parse()
}
