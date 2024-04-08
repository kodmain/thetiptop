package certs

import (
	"crypto/tls"

	"golang.org/x/crypto/acme/autocert"
)

// generateRemoteSignedCert generates a TLS configuration for the provided hosts using remotely signed certificates.
// It leverages the autocert package to manage certificates, including obtaining new certificates from Let's Encrypt,
// renewing them, and storing them in a specified cache directory.
//
// Parameters:
// - hosts: []string A list of hostnames for which the TLS configuration will generate certificates.
//
// Returns:
// - *tls.Config: The TLS configuration with the remotely signed certificates.
func generateRemoteSignedCert(hosts []string) *tls.Config {
	// Initialize a certManager with autocert to manage certificates for the specified hosts
	var certManager = &autocert.Manager{
		Prompt:     autocert.AcceptTOS,               // Automatically accept the terms of service of the CA.
		HostPolicy: autocert.HostWhitelist(hosts...), // Define allowed hosts based on the provided list.
		Cache:      autocert.DirCache("certs"),       // Define a cache directory for storing certificates.
	}

	// Generate and return the TLS configuration using the certManager
	return certManager.TLSConfig()
}
