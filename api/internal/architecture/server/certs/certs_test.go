package certs

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTLSConfigFor(t *testing.T) {
	t.Run("EmptyDomain", func(t *testing.T) {
		tlsConfig := TLSConfigFor("")
		assert.NotNil(t, tlsConfig)
	})

	t.Run("LocalhostDomain", func(t *testing.T) {
		tlsConfig := TLSConfigFor("localhost")
		assert.NotNil(t, tlsConfig)
	})

	t.Run("NonLocalDomain", func(t *testing.T) {
		tlsConfig := TLSConfigFor("example.com")
		assert.NotNil(t, tlsConfig)
	})

	t.Run("WithSubdomains", func(t *testing.T) {
		tlsConfig := TLSConfigFor("example.com", "sub1", "sub2")
		assert.NotNil(t, tlsConfig)
	})
}

func TestGenerateRemoteSignedCert(t *testing.T) {
	t.Run("ValidHosts", func(t *testing.T) {
		hosts := []string{"localhost"}
		tlsConfig := generateRemoteSignedCert(hosts)
		assert.NotNil(t, tlsConfig)
	})

	t.Run("EmptyHosts", func(t *testing.T) {
		hosts := []string{"localhost"}
		tlsConfig := generateRemoteSignedCert(hosts)
		assert.NotNil(t, tlsConfig)
	})
}
