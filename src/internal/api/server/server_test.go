package server

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewServer_Success(t *testing.T) {
	// Arrange

	// Act
	server := NewServer("../../../../src/internal/config/config.yml", "../../../../src/internal/config/kube.config")

	// Assert
	assert.NotNil(t, server)
	assert.NotNil(t, server.KubernetesClient)
}
