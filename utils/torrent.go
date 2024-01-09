package utils

import (
	"fmt"
	"github.com/anacrolix/torrent"
	"github.com/anacrolix/torrent/metainfo"
	"io/ioutil"
	"net/http"
	"os"
)

type TorrentDownload struct {
	Client        *torrent.Client
	Torrent       *torrent.Torrent
	Meta          *metainfo.MetaInfo
	StopRequested bool
	Size          int64 // Size of the downloaded file
}

type TorrentManager struct {
	Downloads    map[string]*TorrentDownload
	DownloadPath string
}

func NewTorrentManager(path string) *TorrentManager {
	return &TorrentManager{
		Downloads:    make(map[string]*TorrentDownload),
		DownloadPath: path,
	}
}

func (m *TorrentManager) AddTorrent(torrentURL string) error {
	// Download the torrent file
	resp, err := http.Get(torrentURL)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return fmt.Errorf("failed to download torrent file: %s", resp.Status)
	}

	// Read the torrent file
	torrentData, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	// Write to a temporary file
	tmpFile, err := ioutil.TempFile("", "*.torrent")
	if err != nil {
		return err
	}
	defer os.Remove(tmpFile.Name()) // Clean up

	if _, err = tmpFile.Write(torrentData); err != nil {
		return err
	}
	if err = tmpFile.Close(); err != nil {
		return err
	}

	cfg := torrent.NewDefaultClientConfig()
	cfg.DataDir = m.DownloadPath

	client, err := torrent.NewClient(cfg)
	if err != nil {
		return err
	}

	tor, err := client.AddTorrentFromFile(tmpFile.Name())
	if err != nil {
		return err
	}
	<-tor.GotInfo()

	// Get the size of the torrent (can be updated later if size changes)
	size := tor.Info().TotalLength()

	download := &TorrentDownload{
		Client:  client,
		Torrent: tor,
		Size:    size,
	}

	m.Downloads[torrentURL] = download
	tor.DownloadAll()
	//m.manageCache() // Call function to manage cache size
	return nil
}

func (m *TorrentManager) CheckProgress() map[string]float64 {
	progress := make(map[string]float64)
	for key, download := range m.Downloads {
		bytesCompleted := download.Torrent.BytesCompleted()
		totalBytes := download.Torrent.Info().TotalLength()
		progress[key] = (float64(bytesCompleted) / float64(totalBytes)) * 100
	}
	return progress
}

func (m *TorrentManager) StopDownload(magnetURI string) {
	if download, ok := m.Downloads[magnetURI]; ok {
		download.StopRequested = true
		download.Torrent.Drop()
		download.Client.Close()
		delete(m.Downloads, magnetURI)
	}
}
