package config

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestNewConfig tests the NewConfig function
func TestNewConfig(t *testing.T) {
	// Create a temporary configuration file
	configContent := `
application:
  grpc-port: "9090"
  http-port: "8080"
  version: "1.0.0"
`
	tmpFile, err := os.CreateTemp("", "config-*.yaml")
	require.NoError(t, err, "failed to create temp file")
	defer os.Remove(tmpFile.Name())

	_, err = tmpFile.Write([]byte(configContent))
	require.NoError(t, err, "failed to write to temp file")
	err = tmpFile.Close()
	require.NoError(t, err, "failed to close temp file")

	// Test NewConfig function
	config, err := NewConfig(tmpFile.Name())
	require.NoError(t, err, "expected no error")

	assert.Equal(t, "9090", config.Application.GRPCPort, "expected grpc-port to be '9090'")
	assert.Equal(t, "8080", config.Application.HTTPPort, "expected http-port to be '8080'")
	assert.Equal(t, "1.0.0", config.Application.Version, "expected version to be '1.0.0'")
}

// TestValidateConfigPath tests the ValidateConfigPath function
func TestValidateConfigPath(t *testing.T) {
	tmpFile, err := os.CreateTemp("", "config-*.yaml")
	require.NoError(t, err, "failed to create temp file")
	defer os.Remove(tmpFile.Name())

	// Test with valid file path
	err = ValidateConfigPath(tmpFile.Name())
	assert.NoError(t, err, "expected no error")

	// Test with directory path
	err = ValidateConfigPath(os.TempDir())
	assert.Error(t, err, "expected error")
}
