package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Mocking the server interface
type MockServer struct {
	mock.Mock
}

func (m *MockServer) Start() error {
	args := m.Called()
	return args.Error(0)
}

// Mocking the NewServer function
var NewServer = func(configPath string, kubeConfigPath string) *MockServer {
	return new(MockServer)
}

// Mocking the AddMiddlewares function
var AddMiddlewares = func(srv *MockServer) {
	srv.On("Start").Return(nil)
}

func TestMainFunction(t *testing.T) {
	// Arrange
	mockServer := new(MockServer)
	NewServer = func(configPath string, kubeConfigPath string) *MockServer {
		return mockServer
	}
	AddMiddlewares = func(srv *MockServer) {
		// Do nothing
	}

	// Act
	srv := NewServer("src/internal/config/config.yml", "src/internal/config/kube.config")
	AddMiddlewares(srv)

	// Assert
	assert.NotNil(t, srv)
	mockServer.AssertExpectations(t)
}
