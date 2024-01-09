package utils

import "os"

// HomeDirGetter is the interface for getting the user's home directory.
type HomeDirGetter interface {
	UserHomeDir() (string, error)
}

// RealHomeDirGetter is the real implementation that uses os.UserHomeDir.
type RealHomeDirGetter struct{}

func (RealHomeDirGetter) UserHomeDir() (string, error) {
	return os.UserHomeDir()
}

// MockHomeDirGetter is the mock implementation for testing.
type MockHomeDirGetter struct {
	HomeDir string
}

func (m MockHomeDirGetter) UserHomeDir() (string, error) {
	return m.HomeDir, nil
}
