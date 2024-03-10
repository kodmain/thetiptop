// Package config provides functions/var/const for loading and accessing configuration settings for the application.
package config

import (
	"github.com/spf13/cobra"
)

// Helper use Cobra package to create a CLI and give Args gesture
var Helper *cobra.Command = &cobra.Command{
	Use:                   "thethiptop",
	Short:                 "TheTipTop API Server",
	DisableAutoGenTag:     true,
	DisableFlagsInUseLine: true,
}

// Declaration of multiple variables used to configure an HTTP application.
var (
	// EnableHTTPS is a boolean variable that indicates whether HTTPS is enabled or not for the application.
	// If set to true, it means the application uses HTTPS for all HTTP connections.
	// Otherwise, it's set to false and the application uses HTTP.
	EnableHTTPS = false
	// TLSPath is a string variable that specifies the absolute path to the directory where SSL certificates are stored.
	TLSPath = "/etc/ssl/certs"
	// TLSCertFile is a string variable that specifies the filename for the SSL certificate used to establish HTTPS connections.
	TLSCertFile = "/server.crt"
	// TLSKeyFile is a string variable that specifies the filename for the SSL private key used to establish HTTPS connections.
	TLSKeyFile = "/server.key"
)

// Initialize Helper
func init() {
	Helper.PersistentFlags().BoolVarP(&EnableHTTPS, "enable-https", "S", EnableHTTPS, "Enable HTTPS")
	Helper.PersistentFlags().StringVarP(&TLSPath, "tls", "t", TLSPath, "define certificats TLS path")
	Helper.PersistentFlags().StringVarP(&TLSCertFile, "cert", "c", TLSCertFile, "define certificats CERT name")
	Helper.PersistentFlags().StringVarP(&TLSKeyFile, "key", "k", TLSKeyFile, "define certificats KEY name")
}
