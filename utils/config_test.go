package utils

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestNewConfig tests the NewConfig function
func TestNewConfig(t *testing.T) {
	// This test should check if a new configuration is correctly created and read

	// Setup: Create a temporary directory to simulate the user's home directory
	tempDir := t.TempDir()
	mockHomeDirGetter := MockHomeDirGetter{HomeDir: tempDir}

	// Act: Call NewConfig with the mock implementation
	config := NewConfig(mockHomeDirGetter)

	// Assert: Check if the configuration was created correctly
	assert.NotNil(t, config, "Config should not be nil")

	os.RemoveAll(tempDir) // Clean up
	_, err := os.Stat(tempDir)
	assert.True(t, os.IsNotExist(err))
}

// TestSave tests the Save method of AppConfig
func TestSave(t *testing.T) {
	// This test should check if the configuration is correctly saved

	// Setup: Create a temporary directory to simulate the user's home directory
	tempDir := t.TempDir()
	mockHomeDirGetter := MockHomeDirGetter{HomeDir: tempDir}

	// Act: Call NewConfig with the mock implementation
	config := NewConfig(mockHomeDirGetter) // Override the HOME environment variable for the test

	// Act: Call Save to save the configuration
	err := config.Save()

	// Assert: Check if the file was saved without errors
	assert.Nil(t, err, "Save should not return an error")

	os.RemoveAll(tempDir) // Clean up
	_, err = os.Stat(tempDir)
	assert.True(t, os.IsNotExist(err))
}
