package installer

import (
	"archive/zip"
	"changeme/utils"
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/saracen/fastzip"
	"io"
	"io/fs"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

type ZiboInstaller struct {
	torrentManager *utils.TorrentManager
	rss            *utils.Rss
	Config         *utils.Config
}

func NewZibo(homeDirGetter utils.HomeDirGetter, singleton bool) *ZiboInstaller {
	config := utils.GetConfig(homeDirGetter, singleton)
	return &ZiboInstaller{
		torrentManager: utils.NewTorrentManager(config.YazuCachePath),
		rss:            utils.NewRss("https://skymatixva.com/tfiles/feed.xml"),
		Config:         config,
	}
}

func (z *ZiboInstaller) GetCachedVersions(update bool) []string {
	// check if path contains a file that contains name of X-Plane 12 Installer
	var cachedVersions []string
	searchDir := filepath.Join(z.Config.YazuCachePath, "full")
	if update {
		searchDir = filepath.Join(z.Config.YazuCachePath, "patch")
	}
	_ = filepath.Walk(searchDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() && path != searchDir {
			// If the path is a subdirectory, skip it
			return filepath.SkipDir
		}
		if !info.IsDir() {
			if strings.LastIndex(path, "zip") != -1 {

				cachedZipPath := path
				version := strings.ReplaceAll(cachedZipPath, searchDir, "")
				version = strings.ReplaceAll(version, ".zip", "")
				version = strings.ReplaceAll(version, "/", "")
				cachedVersions = append(cachedVersions, version)
			}
		}
		return nil
	})

	return cachedVersions
}

func (z *ZiboInstaller) Update(installation utils.ZiboInstallation) {
	patchItems := *z.rss.GetPatchInstallItems()
	patchItem := patchItems[len(patchItems)-1]
	download := z.torrentManager.Downloads[patchItem.Link]
	files := download.Torrent.Files()
	file := files[0]
	_ = os.Rename(filepath.Join(z.Config.YazuCachePath, "patch", file.Path()), filepath.Join(z.Config.YazuCachePath, "patch", patchItem.Version+".zip"))
	z.unzip(filepath.Join(z.Config.YazuCachePath, "patch", patchItem.Version+".zip"), installation.Path)
}

func (z *ZiboInstaller) Backup(installation utils.ZiboInstallation) (string, error) {
	if installation.Path == "" {
		return "", fmt.Errorf("installation path is empty")
	}
	// Check if the directory exists, create it if not
	backupDir := filepath.Join(z.Config.YazuCachePath, "backup")
	if _, err := os.Stat(backupDir); os.IsNotExist(err) {
		err := os.MkdirAll(backupDir, 0755) // Creates the directory with read, write, and execute permissions for the user
		if err != nil {
			log.Fatalf("Error creating Config directory: %s", err)
		}
	}
	// current epoch time
	epochTimeStr := time.Now().Format("2006-01-02_15-04-05")
	zipFilePath := filepath.Join(backupDir, installation.Version+"-"+epochTimeStr+".zip")
	// Create archive file
	w, err := os.Create(zipFilePath)
	if err != nil {
		return "", err
	}
	defer w.Close()

	// Create new Archiver
	a, err := fastzip.NewArchiver(w, installation.Path)
	if err != nil {
		return "", err
	}
	defer a.Close()

	files := make(map[string]os.FileInfo)
	err = filepath.Walk(installation.Path, func(pathname string, info os.FileInfo, err error) error {
		files[pathname] = info
		return nil
	})

	// Archive
	if err = a.Archive(context.Background(), files); err != nil {
		return "", err
	}

	return zipFilePath, nil
}

func (z *ZiboInstaller) Restore(installation utils.ZiboInstallation) error {
	if runtimeOS := runtime.GOOS; runtimeOS == "darwin" {
		// run shell cmd
		script := fmt.Sprintf("do shell script \"sudo rm -rf '%s'\" with administrator privileges", installation.Path)
		log.Println(script)
		std, err := exec.Command("osascript", "-e", script).CombinedOutput()
		if err != nil {
			log.Println(err)
			return err
		}
		log.Println(string(std))
	} else {
		_ = os.RemoveAll(installation.Path)
	}

	backupZip := filepath.Join(z.Config.YazuCachePath, "backup", installation.BackupVersion+".zip")
	destination := installation.Path
	z.unzip(backupZip, destination)
	return nil
}

func (z *ZiboInstaller) unzip(src, dst string, fresh ...bool) {
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

func (z *ZiboInstaller) GetLastBackupVersion() string {
	backupVersion := "N/A"
	backupDir := filepath.Join(z.Config.YazuCachePath, "backup")
	_ = filepath.Walk(backupDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() && path != backupDir {
			// If the path is a subdirectory, skip it
			return filepath.SkipDir
		}
		if !info.IsDir() {
			if strings.LastIndex(path, ".zip") != -1 {
				backupPath := path
				backupVersion = strings.ReplaceAll(backupPath, z.Config.YazuCachePath, "")
				backupVersion = strings.ReplaceAll(backupVersion, ".zip", "")
				backupVersion = strings.ReplaceAll(backupVersion, "/backup/", "")
			}
		}
		return nil
	})

	return backupVersion
}

func (z *ZiboInstaller) RemoveOldInstalls(installation utils.ZiboInstallation) {
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

func (z *ZiboInstaller) DownloadZibo(fullInstall bool) bool {
	var installItem utils.Item
	if fullInstall {
		fullInstallItems := *z.rss.GetFullInstallItems()
		installItem = fullInstallItems[0]
	} else {
		patchedItems := *z.rss.GetPatchInstallItems()
		installItem = patchedItems[len(patchedItems)-1]
	}
	cached := false
	isDownloading := false
	cachedVersions := z.GetCachedVersions(false)

	for _, cachedVersion := range cachedVersions {
		if cachedVersion == installItem.Version && installItem.Version != "" {
			log.Printf("Found cached version %s", cachedVersion)
			cached = true
			break
		}
	}

	if !cached {
		log.Printf("Downloading %s, from: %s", installItem.Version, installItem.Link)
		subPath := "full"
		if !fullInstall {
			subPath = "patch/"
		}
		err := z.torrentManager.AddTorrent(installItem.Link, subPath)
		if err != nil {
			log.Fatalf("Error downloading torrent: %s", err)
		}
		isDownloading = true
	}

	return isDownloading
}

func (z *ZiboInstaller) Install(installation utils.ZiboInstallation) {
	fullInstallItems := *z.rss.GetFullInstallItems()
	fullInstallItem := fullInstallItems[0]
	download := z.torrentManager.Downloads[fullInstallItem.Link]
	files := download.Torrent.Files()
	file := files[0]
	_ = os.Rename(filepath.Join(z.Config.YazuCachePath, file.Path()), filepath.Join(z.Config.YazuCachePath, fullInstallItem.Version+".zip"))
	z.unzip(filepath.Join(z.Config.YazuCachePath, fullInstallItem.Version+".zip"), installation.Path)
}

func (z *ZiboInstaller) GetDownloadProgress(update bool) float64 {
	var link string
	if !update {
		fullItems := *z.rss.GetFullInstallItems()
		fullItem := fullItems[0]
		link = fullItem.Link
	} else {
		patchItems := *z.rss.GetPatchInstallItems()
		patchItem := patchItems[len(patchItems)-1]
		link = patchItem.Link
	}

	progress := z.torrentManager.CheckProgress()
	return progress[link]
}

func (z *ZiboInstaller) FindInstallationDetails() utils.ZiboInstallation {
	var foundPath, version string
	_ = filepath.Walk(filepath.Join(z.Config.XPlanePath, "aircraft"), func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err // prevent panic by handling failure accessing a path
		}
		if info.IsDir() && info.Name() == "zibomod" {
			foundPath = path
			return filepath.SkipDir // folder found, skip the rest of this directory
		}
		return nil
	})
	if foundPath != "" {
		foundPath = filepath.Join(foundPath, "../", "../")
		versionFilePath := filepath.Join(foundPath, "version.txt")

		data, err := os.ReadFile(versionFilePath)
		if err != nil {
			log.Fatalf("Failed to read file: %v", err)
		}
		version = string(data)
	}
	return utils.ZiboInstallation{
		Path:          foundPath,
		Version:       version,
		RemoteVersion: z.rss.GetLatestVersion(),
		BackupVersion: z.GetLastBackupVersion(),
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
