package installer

import (
	"changeme/utils"
	"os"
	"path"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestNewZibo tests the NewZibo function
func TestNewZibo(t *testing.T) {
	// This test should check if a new configuration is correctly created and read

	// Setup: Create a temporary directory to simulate the user's home directory
	tempDir := t.TempDir()
	mockHomeDirGetter := utils.MockHomeDirGetter{HomeDir: tempDir}

	// Act: Call NewZibo with the mock implementation
	config := NewZibo(mockHomeDirGetter)

	// Assert: Check if the configuration was created correctly
	assert.NotNil(t, config, "ZiboInstaller should not be nil")

	os.RemoveAll(tempDir) // Clean up
	_, err := os.Stat(tempDir)
	assert.True(t, os.IsNotExist(err))
}

// TestGetCachedVersions
func TestGetCachedVersions(t *testing.T) {
	t.Run("TestGetCachedVersions - full", func(t *testing.T) {
		// Setup: Create a temporary directory to simulate the user's home directory
		tempDir := t.TempDir()
		mockHomeDirGetter := utils.MockHomeDirGetter{HomeDir: tempDir}

		// Act: Call NewZibo with the mock implementation
		config := NewZibo(mockHomeDirGetter) // Override the HOME environment variable for the test

		// create a file in the cache directory
		cacheDir := tempDir + "/.yazu/cache"
		err := os.MkdirAll(cacheDir, 0755)
		assert.Nil(t, err, "Error creating cache directory")

		cacheFile := cacheDir + "/4.00.1.zip"
		_, err = os.Create(cacheFile)
		assert.Nil(t, err, "Error creating cache file")

		cacheFile = cacheDir + "/4.00.3.zip.patch"
		_, err = os.Create(cacheFile)
		assert.Nil(t, err, "Error creating cache file")

		versions := config.GetCachedVersions(false)
		assert.Equal(t, 1, len(versions))

		os.RemoveAll(tempDir) // Clean up
		_, err = os.Stat(tempDir)
		assert.True(t, os.IsNotExist(err))
	})
	t.Run("TestGetCachedVersions - full - not exist", func(t *testing.T) {
		// Setup: Create a temporary directory to simulate the user's home directory
		tempDir := t.TempDir()
		mockHomeDirGetter := utils.MockHomeDirGetter{HomeDir: tempDir}

		// Act: Call NewZibo with the mock implementation
		config := NewZibo(mockHomeDirGetter) // Override the HOME environment variable for the test

		// create a file in the cache directory
		cacheDir := tempDir + "/.yazu/cache"
		err := os.MkdirAll(cacheDir, 0755)
		assert.Nil(t, err, "Error creating cache directory")

		cacheFile := cacheDir + "/4.00.1.zip.patch"
		_, err = os.Create(cacheFile)
		assert.Nil(t, err, "Error creating cache file")

		cacheFile = cacheDir + "/4.00.3.zip.patch"
		_, err = os.Create(cacheFile)
		assert.Nil(t, err, "Error creating cache file")

		versions := config.GetCachedVersions(false)
		assert.Equal(t, 0, len(versions))

		os.RemoveAll(tempDir) // Clean up
		_, err = os.Stat(tempDir)
		assert.True(t, os.IsNotExist(err))
	})
	t.Run("TestGetCachedVersions - ptch", func(t *testing.T) {
		// Setup: Create a temporary directory to simulate the user's home directory
		tempDir := t.TempDir()
		mockHomeDirGetter := utils.MockHomeDirGetter{HomeDir: tempDir}

		// Act: Call NewZibo with the mock implementation
		config := NewZibo(mockHomeDirGetter) // Override the HOME environment variable for the test

		// create a file in the cache directory
		cacheDir := tempDir + "/.yazu/cache"
		err := os.MkdirAll(cacheDir, 0755)
		assert.Nil(t, err, "Error creating cache directory")

		cacheFile := cacheDir + "/4.00.1.zip"
		_, err = os.Create(cacheFile)
		assert.Nil(t, err, "Error creating cache file")

		cacheFile = cacheDir + "/4.00.3.zip.patch"
		_, err = os.Create(cacheFile)
		assert.Nil(t, err, "Error creating cache file")

		versions := config.GetCachedVersions(true)
		assert.Equal(t, 1, len(versions))
		assert.Equal(t, "4.00.3", versions[0])

		os.RemoveAll(tempDir) // Clean up
		_, err = os.Stat(tempDir)
		assert.True(t, os.IsNotExist(err))
	})
	t.Run("TestGetCachedVersions - ptch - not exist", func(t *testing.T) {
		// Setup: Create a temporary directory to simulate the user's home directory
		tempDir := t.TempDir()
		mockHomeDirGetter := utils.MockHomeDirGetter{HomeDir: tempDir}

		// Act: Call NewZibo with the mock implementation
		config := NewZibo(mockHomeDirGetter) // Override the HOME environment variable for the test

		// create a file in the cache directory
		cacheDir := tempDir + "/.yazu/cache"
		err := os.MkdirAll(cacheDir, 0755)
		assert.Nil(t, err, "Error creating cache directory")

		cacheFile := cacheDir + "/4.00.1.zip"
		_, err = os.Create(cacheFile)
		assert.Nil(t, err, "Error creating cache file")

		cacheFile = cacheDir + "/4.00.3.zip"
		_, err = os.Create(cacheFile)
		assert.Nil(t, err, "Error creating cache file")

		versions := config.GetCachedVersions(true)
		assert.Equal(t, 0, len(versions))

		os.RemoveAll(tempDir) // Clean up
		_, err = os.Stat(tempDir)
		assert.True(t, os.IsNotExist(err))
	})
}

// TestInstall
func TestInstall(t *testing.T) {
	t.Run("TestInstall - full", func(t *testing.T) {
		// Setup: Create a temporary directory to simulate the user's home directory
		tempDir := t.TempDir()
		mockHomeDirGetter := utils.MockHomeDirGetter{HomeDir: tempDir}

		// Act: Call NewZibo with the mock implementation
		config := NewZibo(mockHomeDirGetter) // Override the HOME environment variable for the test

		// create a file in the cache directory
		cacheDir := tempDir + "/.yazu/cache"
		err := os.MkdirAll(cacheDir, 0755)
		assert.Nil(t, err, "Error creating cache directory")

		cacheFile := cacheDir + "/4.00.1.zip"
		_, err = os.Create(cacheFile)
		assert.Nil(t, err, "Error creating cache file")

		cacheFile = cacheDir + "/4.00.3.zip.patch"
		_, err = os.Create(cacheFile)
		assert.Nil(t, err, "Error creating cache file")

		config.Install(utils.ZiboInstallation{
			Path:          path.Join(tempDir, "Aircraft", "Zibo B738X"),
			Version:       "",
			RemoteVersion: "4.00.1",
			BackupVersion: "",
		})

		os.RemoveAll(tempDir) // Clean up
		_, err = os.Stat(tempDir)
		assert.True(t, os.IsNotExist(err))
	})
	//t.Run("TestInstall - full - not exist", func(t *testing.T) {
	//	// Setup: Create a temporary directory to simulate the user's home directory
	//	tempDir := t.TempDir()
	//	mockHomeDirGetter := MockHomeDirGetter{HomeDir: tempDir}
	//
	//	// Act: Call NewZibo with the mock implementation
	//	Config := NewZibo(mockHomeDirGetter) // Override the HOME environment variable for the test
	//
	//	// create a file in the cache directory
	//	cacheDir := tempDir + "/.yazu/cache"
	//	err := os.MkdirAll(cacheDir, 0755)
	//	assert.Nil(t, err, "Error creating cache directory")
	//
	//	cacheFile := cacheDir + "/4.00.1.zip.patch"
	//	_, err = os.Create(cacheFile)
	//	assert.Nil(t, err, "Error creating cache file")
	//
	//	cacheFile = cacheDir + "/4.00.3.zip.patch"
	//	_, err = os.Create(cacheFile)
	//	assert.Nil(t, err, "Error creating cache file")
	//}
}
