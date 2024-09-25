package config_test

import (
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/kodmain/thetiptop/api/config"
	"github.com/stretchr/testify/assert"
)

func TestLoad(t *testing.T) {
	assert.Nil(t, config.Get("mail", nil))
	err := config.Load(aws.String(""))
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "path is required")

	err = config.Load(aws.String("cnf.yml"))
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "open cnf.yml: no such file or directory")

	err = config.Load(aws.String("../config.test.yml"))
	assert.Nil(t, err)

	err = config.Load(aws.String("s3://config.kodmain/config.yml"))
	assert.Error(t, err) // This test will fail because no credentials are provided
}

func TestGetString(t *testing.T) {
	config.Load(aws.String("../config.test.yml"))
	assert.Equal(t, "default", config.GetString("providers.databases", "default"))
	assert.Equal(t, "default", config.GetString("services.client.database", "default"))
}

func TestGetInt(t *testing.T) {
	config.Load(aws.String("../config.test.yml"))
	assert.Equal(t, 3, config.GetInt("providers.databases", 3))
	assert.Equal(t, 1025, config.GetInt("providers.mails.default.port", 3))
}

func TestAll(t *testing.T) {
	config.Load(aws.String("../config.test.yml"))

	// Services
	assert.Equal(t, "default", config.Get("services.client.database", "default-value"))
	assert.Equal(t, "default", config.Get("services.client.mail", "default-value"))
	assert.Equal(t, "default", config.Get("services.employee.database", "default-value"))
	assert.Equal(t, "default", config.Get("services.employee.mail", "default-value"))

	// Providers - mails
	assert.Equal(t, "secret", config.Get("providers.mails.default.username", "default-value"))
	assert.Equal(t, "secret", config.Get("providers.mails.default.password", "default-value"))
	assert.Equal(t, "localhost", config.Get("providers.mails.default.host", "default-value"))
	assert.Equal(t, 1025, config.GetInt("providers.mails.default.port", 0))
	assert.Equal(t, "Whoami", config.Get("providers.mails.default.expeditor", "default-value"))
	assert.Equal(t, "whoami@localhost", config.Get("providers.mails.default.from", "default-value"))

	// Providers - databases
	assert.Equal(t, "sqlite", config.Get("providers.databases.file.protocol", "default-value"))
	assert.Contains(t, config.Get("providers.databases.file.dbname", "default-value"), "db.sqlite")
	assert.Equal(t, true, config.Get("providers.databases.file.logger", false))

	assert.Equal(t, "sqlite", config.Get("providers.databases.default.protocol", "default-value"))
	assert.Equal(t, ":memory:", config.Get("providers.databases.default.dbname", "default-value"))
	assert.Equal(t, true, config.Get("providers.databases.default.logger", false))

	// Security - validation
	assert.Equal(t, "30m", config.Get("security.validation.expire", "default-value"))

	// Security - jwt
	assert.Equal(t, "Europe/Paris", config.Get("security.jwt.tz", "default-value"))
	assert.Equal(t, "secret", config.Get("security.jwt.secret", "default-value"))
	assert.Equal(t, 15, config.GetInt("security.jwt.expire", 0))
	assert.Equal(t, 30, config.GetInt("security.jwt.refresh", 0))
}

func TestGet(t *testing.T) {
	// Reset config before running tests
	config.Reset()

	// Test retrieving non-existent value before loading any config
	assert.Nil(t, config.Get("mail", nil))

	// Load valid config file for testing
	config.Load(aws.String("../config.test.yml"))

	// Test retrieving non-existent key with and without default value
	assert.Nil(t, config.Get("toto", nil))
	assert.Equal(t, "toto", config.Get("toto", "toto"))

	// Test retrieving valid existing key
	assert.NotNil(t, config.Get("security", nil))
	assert.NotNil(t, config.Get("security.jwt", nil))

	// Test retrieving nested map fields
	assert.NotNil(t, config.Get("services", nil))
	assert.NotNil(t, config.Get("providers", nil))

	// Test retrieving map elements by key
	assert.NotNil(t, config.Get("providers.mails", nil))
	assert.NotNil(t, config.Get("providers.databases", nil))

	// Test retrieving nested values inside maps
	assert.NotNil(t, config.Get("providers.mails.default", nil))
	assert.NotNil(t, config.Get("providers.databases.default", nil))

	// Additional test: ensure default value is returned when accessing nil pointer
	assert.Equal(t, "defaultMail", config.Get("non.existent.mail", "defaultMail"))

	// Additional test: map key exists but value is nil, should return default
	assert.Nil(t, config.Get("services.nilService", nil))
	assert.Equal(t, "defaultService", config.Get("services.nilService", "defaultService"))

	// Test map access with invalid map keys and ensure default value is returned
	assert.Equal(t, "defaultKey", config.Get("providers.mails.invalidKey", "defaultKey"))

	// Additional test: nested field retrieval with invalid map path
	assert.Equal(t, "defaultJWT", config.Get("security.invalid.jwt", "defaultJWT"))
}
