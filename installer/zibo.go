package installer

import (
	"archive/zip"
	"changeme/utils"
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/saracen/fastzip"
	"gopkg.in/yaml.v3"
	"io"
	"io/fs"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"runtime"
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
	torrentManager *utils.TorrentManager
	mu             sync.Mutex
	rss            *utils.Rss
}

func NewConfig(homeDirGetter utils.HomeDirGetter) *Config {
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

	config.torrentManager = utils.NewTorrentManager(config.YazuCachePath)
	config.rss = utils.NewRss("https://skymatixva.com/tfiles/feed.xml")
	return config
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

func (c *Config) GetCachedVersions(update bool) []string {
	// check if path contains a file that contains name of X-Plane 12 Installer
	cachedVersions := []string{}
	_ = filepath.Walk(c.YazuCachePath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() && path != c.YazuCachePath {
			// If the path is a subdirectory, skip it
			return filepath.SkipDir
		}
		if !info.IsDir() {
			if strings.LastIndex(path, "zip") != -1 && strings.LastIndex(path, "backup") == -1 {
				if !update {
					if strings.LastIndex(path, "patch") != -1 {
						return nil
					}
				} else {
					if strings.LastIndex(path, "patch") == -1 {
						return nil
					}
				}

				cachedZipPath := path
				version := strings.ReplaceAll(cachedZipPath, c.YazuCachePath, "")
				version = strings.ReplaceAll(version, ".zip", "")
				version = strings.ReplaceAll(version, "/", "")
				if update {
					version = strings.ReplaceAll(version, ".patch", "")
				}
				cachedVersions = append(cachedVersions, version)
			}
		}
		return nil
	})

	return cachedVersions
}

func (c *Config) Update(installation ZiboInstallation) {
	patchItems := *c.rss.GetPatchItems()
	patchItem := patchItems[len(patchItems)-1]
	download := c.torrentManager.Downloads[patchItem.Link]
	files := download.Torrent.Files()
	file := files[0]
	_ = os.Rename(filepath.Join(c.YazuCachePath, file.Path()), filepath.Join(c.YazuCachePath, patchItem.Version+".zip.patch"))
	c.unzip(filepath.Join(c.YazuCachePath, patchItem.Version+".zip.patch"), installation.Path)
}

func (c *Config) Backup(installation ZiboInstallation) bool {
	// Check if the directory exists, create it if not
	if _, err := os.Stat(c.YazuCachePath); os.IsNotExist(err) {
		err := os.MkdirAll(c.YazuCachePath, 0755) // Creates the directory with read, write, and execute permissions for the user
		if err != nil {
			log.Fatalf("Error creating config directory: %s", err)
		}
	}

	zipFilePath := filepath.Join(c.YazuCachePath, installation.Version+"-backup.zip")
	// Create archive file
	w, err := os.Create(zipFilePath)
	if err != nil {
		panic(err)
	}
	defer w.Close()

	// Create new Archiver
	a, err := fastzip.NewArchiver(w, installation.Path)
	if err != nil {
		panic(err)
	}
	defer a.Close()

	files := make(map[string]os.FileInfo)
	err = filepath.Walk(installation.Path, func(pathname string, info os.FileInfo, err error) error {
		files[pathname] = info
		return nil
	})

	// Archive
	if err = a.Archive(context.Background(), files); err != nil {
		panic(err)
	}

	return true
}

func (c *Config) Restore(installation ZiboInstallation) bool {
	if runtimeOS := runtime.GOOS; runtimeOS == "darwin" {
		// run shell cmd
		script := fmt.Sprintf("do shell script \"sudo rm -rf '%s'\" with administrator privileges", installation.Path)
		log.Println(script)
		std, err := exec.Command("osascript", "-e", script).CombinedOutput()
		if err != nil {
			log.Println(err)
		}
		log.Println(string(std))
	} else {
		_ = os.RemoveAll(installation.Path)
	}

	backupZip := filepath.Join(c.YazuCachePath, installation.BackupVersion+"-backup.zip")
	destination := installation.Path
	c.unzip(backupZip, destination)
	return true
}

func (c *Config) unzip(src, dst string, fresh ...bool) {
	// create a tmp directory
	tmpDir := os.TempDir()
	uuid := uuid.New().String()
	_ = os.MkdirAll(filepath.Join(tmpDir, uuid), 0700)

	// Create new extractor
	log.Println("Extracting archive..." + src)

	r, err := zip.OpenReader(src)
	if err != nil {
		log.Fatalf("Error opening zip file: %s", err)
	}
	defer r.Close()

	for _, f := range r.File {
		rc, err := f.Open()
		if err != nil {
			log.Fatalf("Error opening file in zip: %s", err)
		}
		defer rc.Close()

		path := filepath.Join(tmpDir, uuid, f.Name)
		if f.FileInfo().IsDir() {
			os.MkdirAll(path, 0700)
		} else {
			f, err := os.OpenFile(
				path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
			if err != nil {
				log.Fatalf("Error creating file: %s", err)
			}
			defer f.Close()

			_, err = io.Copy(f, rc)
			if err != nil {
				log.Fatalf("Error copying file contents: %s", err)
			}
		}
	}

	// move files from tmp directory to destination
	if runtimeOS := runtime.GOOS; runtimeOS == "darwin" {
		// run shell cmd
		script := fmt.Sprintf("do shell script \"sudo ditto '%s/%s' '%s';sudo xattr -d -r com.apple.quarantine '%s'\" with administrator privileges", tmpDir, uuid, dst, dst)
		log.Printf("Move files: %s", script)
		std, err := exec.Command("osascript", "-e", script).CombinedOutput()
		if err != nil {
			log.Println(err)
		}
		log.Println(string(std))
	} else {
		_ = ditto(filepath.Join(tmpDir, "B737-800X"), dst)
	}

}

// copyFile copies a single file from src to dst, preserving file permissions.
func copyFile(src, dst string) error {
	sourceFileStat, err := os.Stat(src)
	if err != nil {
		return err
	}

	if !sourceFileStat.Mode().IsRegular() {
		return &fs.PathError{Op: "copy", Path: src, Err: fs.ErrInvalid}
	}

	source, err := os.Open(src)
	if err != nil {
		return err
	}
	defer source.Close()

	destination, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer destination.Close()

	_, err = io.Copy(destination, source)
	if err != nil {
		return err
	}

	return os.Chmod(dst, sourceFileStat.Mode())
}

// ditto mimics the basic behavior of the 'ditto' command for directories.
func ditto(src, dst string) error {
	// Get properties of the source directory
	info, err := os.Stat(src)
	if err != nil {
		return err
	}

	// Create the destination directory
	err = os.MkdirAll(dst, info.Mode())
	if err != nil {
		return err
	}

	return filepath.WalkDir(src, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		targetPath := filepath.Join(dst, path[len(src):])

		if d.IsDir() {
			// Create sub-directories.
			return os.MkdirAll(targetPath, d.Type())
		} else {
			// Copy files.
			return copyFile(path, targetPath)
		}
	})
}

func (c *Config) GetLastBackupVersion() string {
	backupVersion := "N/A"
	_ = filepath.Walk(c.YazuCachePath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() && path != c.YazuCachePath {
			// If the path is a subdirectory, skip it
			return filepath.SkipDir
		}
		if !info.IsDir() {
			if strings.LastIndex(path, "-backup.zip") != -1 {
				backupPath := path
				backupVersion = strings.ReplaceAll(backupPath, c.YazuCachePath, "")
				backupVersion = strings.ReplaceAll(backupVersion, "-backup.zip", "")
				backupVersion = strings.ReplaceAll(backupVersion, "/", "")
			}
		}
		return nil
	})

	return backupVersion
}

func (c *Config) RemoveOldInstalls(installation ZiboInstallation) {
	log.Printf("Removing %s", installation.Path)
	// if os is mac
	if runtimeOS := runtime.GOOS; runtimeOS == "darwin" {
		// run shell cmd
		script := fmt.Sprintf("do shell script \"sudo rm -rf '%s'\" with administrator privileges", installation.Path)
		std, err := exec.Command("osascript", "-e", script).CombinedOutput()
		if err != nil {
			log.Println(err)
		}
		log.Println(string(std))
	} else {
		_ = os.RemoveAll(installation.Path)
	}
}

func (c *Config) DownloadZibo(fullInstall bool) bool {
	var installItem utils.Item
	if fullInstall {
		fullInstallItems := *c.rss.GetFullInstallItems()
		installItem = fullInstallItems[0]
	} else {
		patchedItems := *c.rss.GetPatchItems()
		installItem = patchedItems[len(patchedItems)-1]
	}
	cached := false
	isDownloading := false
	cachedVersions := c.GetCachedVersions(false)

	for _, cachedVersion := range cachedVersions {
		if cachedVersion == installItem.Version && installItem.Version != "" {
			log.Printf("Found cached version %s", cachedVersion)
			cached = true
			break
		}
	}

	if !cached {
		log.Printf("Downloading %s, from: %s", installItem.Version, installItem.Link)
		err := c.torrentManager.AddTorrent(installItem.Link)
		if err != nil {
			log.Fatalf("Error downloading torrent: %s", err)
		}
		isDownloading = true
	}

	return isDownloading
}

func (c *Config) Install(installation ZiboInstallation) {
	fullInstallItems := *c.rss.GetFullInstallItems()
	fullInstallItem := fullInstallItems[0]
	download := c.torrentManager.Downloads[fullInstallItem.Link]
	files := download.Torrent.Files()
	file := files[0]
	_ = os.Rename(filepath.Join(c.YazuCachePath, file.Path()), filepath.Join(c.YazuCachePath, fullInstallItem.Version+".zip"))
	c.unzip(filepath.Join(c.YazuCachePath, fullInstallItem.Version+".zip"), installation.Path)
}

func (c *Config) GetDownloadProgress(update bool) float64 {
	var link string
	if !update {
		fullItems := *c.rss.GetFullInstallItems()
		fullItem := fullItems[0]
		link = fullItem.Link
	} else {
		patchItems := *c.rss.GetPatchItems()
		patchItem := patchItems[len(patchItems)-1]
		link = patchItem.Link
	}

	progress := c.torrentManager.CheckProgress()
	return progress[link]
}
