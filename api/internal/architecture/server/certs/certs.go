package certs

import (
	"crypto/tls"
	"fmt"
)

// TLSConfigFor generates a TLS configuration for the given domain and subdomains.
// If the domain is "localhost", it generates a self-signed certificate; otherwise,
// it obtains a remotely signed certificate. This function ensures that the generated
// TLS configuration is appropriate for the specified domain and subdomains.
//
// Parameters:
// - domain: string The primary domain for the TLS configuration.
// - subs: ...string A variadic list of subdomains associated with the domain.
//
// Returns:
// - *tls.Config: The TLS configuration with the appropriate certificate.
func TLSConfigFor(domain string, subs ...string) *tls.Config {
	// Default to "localhost" if the domain is empty
	if domain == "" {
		domain = "localhost"
	}

	// Generate a list of hostnames including the domain and subdomains
	hosts := generateHosts(domain, subs...)

	// Generate a self-signed certificate for localhost, otherwise get a remotely signed certificate
	if domain == "localhost" || domain == "127.0.0.1" {
		return generateSelfSignedCert(hosts)
	}

	return generateRemoteSignedCert(hosts)
}

// generateHosts creates a list of hostnames from a domain and its subdomains.
// It combines the main domain with each subdomain to form complete hostnames.
//
// Parameters:
// - domain: string The main domain.
// - subs: ...string A variadic list of subdomains.
//
// Returns:
// - []string: A list of combined domain and subdomain hostnames.
func generateHosts(domain string, subs ...string) []string {
	var combinations []string
	combinations = append(combinations, domain)
	for _, sub := range subs {
		// Combine each subdomain with the main domain
		combinations = append(combinations, fmt.Sprintf("%s.%s", sub, domain))
	}

	return combinations
}
