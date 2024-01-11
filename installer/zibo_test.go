package installer

import (
	uuid2 "github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"github/xairline/yazu/utils"
	"os"
	"path"
	"testing"

	"github.com/stretchr/testify/assert"
)

func setup(t *testing.T) (*ZiboInstaller, string) {
	// Setup: Create a temporary directory to simulate the user's home directory
	tempDir := t.TempDir()
	uuid := uuid2.New()
	tempDir = path.Join(tempDir, uuid.String())
	_ = os.MkdirAll(tempDir, 0755)
	mockHomeDirGetter := utils.MockHomeDirGetter{HomeDir: tempDir}

	// Act: Call NewZibo with the mock implementation
	return NewZibo(mockHomeDirGetter, false, logrus.New()), tempDir
}

// TestNewZibo tests the NewZibo function
func TestNewZibo(t *testing.T) {
	config, tempDir := setup(t)
	defer func() {
		_ = os.RemoveAll(tempDir) // Clean up
		_, err := os.Stat(tempDir)
		assert.True(t, os.IsNotExist(err))
	}()

	// Assert: Check if the configuration was created correctly
	assert.NotNil(t, config, "ZiboInstaller should not be nil")
}

func TestOrgZiboLiveries(t *testing.T) {
	config, tempDir := setup(t)
	defer func() {
		_ = os.RemoveAll(tempDir) // Clean up
		_, err := os.Stat(tempDir)
		assert.True(t, os.IsNotExist(err))
	}()
	res := config.GetAvailableLiveries()
	// Assert: Check if the configuration was created correctly
	assert.NotNil(t, res)
}

//
//// TestInstall
//func TestInstall(t *testing.T) {
//	t.Run("TestInstall - full", func(t *testing.T) {
//		config, tempDir := setup(t)
//		defer func() {
//			_ = os.RemoveAll(tempDir) // Clean up
//			_, err := os.Stat(tempDir)
//			assert.True(t, os.IsNotExist(err))
//		}()
//
//		// create a file in the cache directory
//		cacheDir := tempDir + "/.yazu/cache"
//		err := os.MkdirAll(cacheDir, 0755)
//		err = os.MkdirAll(cacheDir+"/full", 0755)
//		err = os.MkdirAll(cacheDir+"/patch", 0755)
//		assert.Nil(t, err, "Error creating cache directory")
//
//		cacheFile := cacheDir + "/full/4.00.1.zip"
//		_, err = os.Create(cacheFile)
//		assert.Nil(t, err, "Error creating cache file")
//
//		cacheFile = cacheDir + "/patch/4.00.3.zip"
//		_, err = os.Create(cacheFile)
//		assert.Nil(t, err, "Error creating cache file")
//
//		config.Install(utils.ZiboInstallation{
//			Path:          path.Join(tempDir, "Aircraft", "Zibo B738X"),
//			Version:       "",
//			RemoteVersion: "4.00.1",
//			BackupVersion: "",
//		})
//
//	})
//	//t.Run("TestInstall - full - not exist", func(t *testing.T) {
//	//	// Setup: Create a temporary directory to simulate the user's home directory
//	//	tempDir := t.TempDir()
//	//	mockHomeDirGeter := MockHomeDirGeter{HomeDir: tempDir}
//	//
//	//	// Act: Call NewZibo with the mock implementation
//	//	Config := NewZibo(mockHomeDirGeter) // Override the HOME environment variable for the test
//	//
//	//	// create a file in the cache directory
//	//	cacheDir := tempDir + "/.yazu/cache"
//	//	err := os.MkdirAll(cacheDir, 0755)
//	//	assert.Nil(t, err, "Error creating cache directory")
//	//
//	//	cacheFile := cacheDir + "/4.00.1.zip.patch"
//	//	_, err = os.Create(cacheFile)
//	//	assert.Nil(t, err, "Error creating cache file")
//	//
//	//	cacheFile = cacheDir + "/4.00.3.zip.patch"
//	//	_, err = os.Create(cacheFile)
//	//	assert.Nil(t, err, "Error creating cache file")
//	//}
//}
//
//// TestBackup
//func TestBackupRestore(t *testing.T) {
//	t.Run("TestBackup - no install", func(t *testing.T) {
//		zibo, tempDir := setup(t)
//		defer func() {
//			_ = os.RemoveAll(tempDir) // Clean up
//			_, err := os.Stat(tempDir)
//			assert.True(t, os.IsNotExist(err))
//		}()
//		zipFilePath, err := zibo.Backup(zibo.FindInstallationDetails())
//		assert.NotNil(t, err)
//		assert.Equal(t, "", zipFilePath)
//	})
//
//	t.Run("TestBackupRestore - install", func(t *testing.T) {
//		zibo, tempDir := setup(t)
//		defer func() {
//			_ = os.RemoveAll(tempDir) // Clean up
//			_, err := os.Stat(tempDir)
//			assert.True(t, os.IsNotExist(err))
//		}()
//		fakePath := path.Join(tempDir, "Aircraft", "Zibo B738X")
//		os.MkdirAll(fakePath, 0755)
//		os.Create(path.Join(fakePath, "test.txt"))
//
//		zipFilePath, err := zibo.Backup(utils.ZiboInstallation{
//			Path:    fakePath,
//			Version: "fake",
//		})
//		assert.Nil(t, err)
//		assert.FileExists(t, zipFilePath)
//		assert.Contains(t, zipFilePath, "backup/fake-")
//		assert.Contains(t, zipFilePath, ".zip")
//
//		details := zibo.FindInstallationDetails()
//		details.Path = fakePath
//		err = zibo.Restore(details)
//		assert.Nil(t, err)
//		assert.FileExists(t, path.Join(fakePath, "test.txt"))
//	})
//}
