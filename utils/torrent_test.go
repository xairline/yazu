package utils

import (
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"os"
	"testing"
	"time"
)

func TestTorrentDownloadAndStop(t *testing.T) {
	// Setup a temporary directory for downloads
	tmpDir, err := ioutil.TempDir("", "torrent_test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}

	// Initialize TorrentManager with the temp directory
	manager := NewTorrentManager("https://skymatixva.com/tfiles/feed.xml", logrus.New())
	manager.DownloadPath = tmpDir

	// Start a torrent download
	magnetURI := "https://skymatixva.com/tfiles/B738X_XP12_4_00_3.zip.torrent"
	err = manager.AddTorrent(magnetURI, "")
	if err != nil {
		t.Fatalf("Failed to add torrent: %v", err)
	}

	// Allow some time for the download to start
	time.Sleep(10 * time.Second)

	// Check if a file is being created in the specified directory
	files, _ := os.ReadDir(tmpDir)
	if len(files) == 0 {
		t.Error("No files found in download directory, expected at least one")
	}

	progress := manager.CheckProgress()
	assert.NotNil(t, progress[magnetURI])

	// Stop the download
	manager.StopDownload(magnetURI)

	os.RemoveAll(tmpDir) // Clean up
	_, err = os.Stat(tmpDir)
	assert.True(t, os.IsNotExist(err))
}
