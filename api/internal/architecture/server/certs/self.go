package certs

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"math/big"
	"time"

	"github.com/kodmain/thetiptop/api/internal/architecture/observability/logger"
)

// generateSelfSignedCert generates a self-signed certificate and returns a tls.Config.
// This function creates an ECDSA private key, generates a unique serial number, and
// then creates an X.509 certificate with this information. The certificate is valid for
// the provided host names and for a duration of one year. The returned tls.Config is configured
// with the generated certificate and security settings.
//
// Parameters:
// - hosts: []string The host names for which the certificate is valid.
//
// Returns:
// - *tls.Config: TLS configuration containing the self-signed certificate.
func generateSelfSignedCert(hosts []string) *tls.Config {
	// Generate an ECDSA private key
	priv, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if logger.Fatal(err) {
		return nil
	}

	// Generate a serial number for the certificate
	serialNumber, err := rand.Int(rand.Reader, new(big.Int).Lsh(big.NewInt(1), 128))
	if logger.Fatal(err) {
		return nil
	}

	// Define the X.509 certificate template
	template := x509.Certificate{
		SerialNumber: serialNumber,
		Subject: pkix.Name{
			Organization: []string{"Kitsune"},
		},
		NotBefore:             time.Now(),
		NotAfter:              time.Now().Add(365 * 24 * time.Hour),
		KeyUsage:              x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		BasicConstraintsValid: true,
		DNSNames:              hosts,
	}

	// Create the certificate from the template and private key
	derBytes, err := x509.CreateCertificate(rand.Reader, &template, &template, &priv.PublicKey, priv)
	if logger.Fatal(err) {
		return nil
	}

	// Create a tls.Certificate object with the generated certificate
	cert := tls.Certificate{
		Certificate: [][]byte{derBytes},
		PrivateKey:  priv,
	}

	// Configure tls.Config with the certificate and security settings
	tlsConfig := &tls.Config{
		MinVersion: tls.VersionTLS12,
		CipherSuites: []uint16{
			tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
		},
		InsecureSkipVerify: true,
		Certificates:       []tls.Certificate{cert},
	}

	return tlsConfig
}
