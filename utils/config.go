package utils

import (
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"log"
	"os"
	"path"
	"path/filepath"
	"strings"
	"sync"
)

type ZiboInstallation struct {
	Path          string `json:"path"`
	Version       string `json:"version"`
	RemoteVersion string `json:"remoteVersion"`
	BackupVersion string `json:"backupVersion"`
}

type Config struct {
	XPlanePath     string `yaml:"xplane_path" json:"XPlanePath"`
	YazuCachePath  string `yaml:"yazu_cache_path" json:"YazuCachePath"`
	torrentManager *TorrentManager
	mu             sync.Mutex
	rss            *Rss
}

var (
	// instance holds the single instance of Config
	instance *Config
	// once is used to initialize the singleton instance once
	once sync.Once
)

// GetConfig returns the singleton instance of Config.
func GetConfig(homeDirGetter HomeDirGetter, singleton bool) *Config {
	if !singleton {
		return newConfig(homeDirGetter)
	}
	once.Do(func() {
		instance = newConfig(homeDirGetter)
	})
	return instance
}

func newConfig(homeDirGetter HomeDirGetter) *Config {
	// check if config file exists
	homeDir, err := homeDirGetter.UserHomeDir()
	if err != nil {
		log.Fatalf("failed to get user home directory: %v", err)
	}

	configDir := filepath.Join(homeDir, ".yazu")
	configFilePath := filepath.Join(configDir, "config.yaml")

	// Check if the directory exists, create it if not
	if _, err := os.Stat(configDir); os.IsNotExist(err) {
		err := os.MkdirAll(configDir, 0755) // Creates the directory with read, write, and execute permissions for the user
		if err != nil {
			log.Fatalf("Error creating config directory: %s", err)
		}
	}

	// Check if the config file exists
	if _, err := os.Stat(configFilePath); os.IsNotExist(err) {
		// File does not exist, create a default config
		defaultConfig := &Config{
			YazuCachePath: path.Join(homeDir, ".yazu", "cache"),
		}
		data, err := yaml.Marshal(defaultConfig)
		if err != nil {
			log.Fatalf("Error marshaling default config: %s", err)
		}
		if err := ioutil.WriteFile(configFilePath, data, 0644); err != nil {
			log.Fatalf("Error writing default config file: %s", err)
		}
	}
	// File exists, read the config file
	data, err := os.ReadFile(configFilePath)
	if err != nil {
		log.Fatalf("Error reading config file: %s", err)
	}

	// Unmarshal the config data into an Config instance
	config := &Config{}
	if err := yaml.Unmarshal(data, config); err != nil {
		log.Fatalf("Error unmarshaling config data: %s", err)
	}

	config.torrentManager = NewTorrentManager(config.YazuCachePath)
	config.rss = NewRss("https://skymatixva.com/tfiles/feed.xml")
	return config
}

func (c *Config) CheckXPlanePath(dirPath string) bool {
	// check if path contains a file that contains name of X-Plane 12 Installer
	found := false
	_ = filepath.Walk(dirPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() && path != dirPath {
			// If the path is a subdirectory, skip it
			return filepath.SkipDir
		}
		if !info.IsDir() {
			if strings.Contains(path, "Log.txt") {
				// store path in home directory
				c.XPlanePath = filepath.Dir(path)
				_ = c.Save()
				found = true
			}
		}
		return nil
	})

	return found
}

func (c *Config) Save() error {
	// Lock the config mutex to prevent concurrent writes
	c.mu.Lock()
	defer c.mu.Unlock()

	// Marshal the config instance into JSON
	data, err := yaml.Marshal(c)
	if err != nil {
		return err
	}

	// Write the config file
	configFilePath, err := os.UserHomeDir()
	if err != nil {
		log.Fatalf("failed to get user home directory: %v", err)
	}
	configFilePath += "/.yazu/config.yaml"
	if err := os.WriteFile(configFilePath, data, 0644); err != nil {
		return err
	}

	return nil
}
