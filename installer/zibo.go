package installer

import (
	"archive/zip"
	"context"
	"encoding/base64"
	"fmt"
	"github.com/google/uuid"
	"github.com/pkg/xattr"
	"github.com/saracen/fastzip"
	"github.com/sirupsen/logrus"
	"github/xairline/yazu/utils"
	"io"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

type ZiboInstaller struct {
	TorrentManager *utils.TorrentManager
	rss            *utils.Rss
	Config         *utils.Config
	log            *logrus.Logger
}

type ZiboBackup struct {
	BackupPath string `json:"backupPath"`
	Version    string `json:"version"`
	Date       string `json:"date"`
	Size       int    `json:"size"`
}

type InstalledLivery struct {
	Name string `json:"name"`
	Path string `json:"path"`
	Icon string `json:"icon"`
}

type AvailableLivery struct {
	Name   string `json:"name"`
	Url    string `json:"url"`
	Source string `json:"source"`
	Icon   string `json:"icon"`
}

func NewZibo(homeDirGetter utils.HomeDirGetter, singleton bool, log *logrus.Logger) *ZiboInstaller {
	config := utils.GetConfig(homeDirGetter, singleton, log)
	return &ZiboInstaller{
		TorrentManager: utils.NewTorrentManager(config.YazuCachePath, log),
		rss:            utils.NewRss("https://skymatixva.com/tfiles/feed.xml", log),
		Config:         config,
		log:            log,
	}
}

func (z *ZiboInstaller) Update(installation utils.ZiboInstallation, zipFilePath string) {
	z.unzip(zipFilePath, installation.Path, false)
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
			z.log.Errorf("Error creating Config directory: %s", err)
			return "", err
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

func (z *ZiboInstaller) Restore(installation utils.ZiboInstallation, backupPath string) error {
	z.RemoveOldInstalls(installation)

	backupZip := filepath.Join(z.Config.YazuCachePath, "backup", installation.BackupVersion+".zip")
	if backupPath != "" {
		backupZip = backupPath
	}
	destination := installation.Path
	z.unzip(backupZip, destination, false)
	return nil
}

func (z *ZiboInstaller) unzip(src, dst string, fresh bool) {
	// create a tmp directory
	tmpDir := os.TempDir()
	uuid := uuid.New().String()
	if fresh {
		uuid = uuid + "/B737-800X"
		dst = filepath.Join(dst, "..")
	}
	tmpUnzipDir := filepath.Join(tmpDir, uuid)
	_ = os.MkdirAll(tmpUnzipDir, 0700)

	// Create new extractor
	z.log.Infof("Extracting archive..." + src)

	r, err := zip.OpenReader(src)
	if err != nil {
		z.log.Errorf("Error opening zip file: %s", err)
	}
	defer r.Close()

	for _, f := range r.File {
		rc, err := f.Open()
		if err != nil {
			z.log.Errorf("Error opening file in zip: %s", err)
		}
		defer rc.Close()

		path := filepath.Join(tmpDir, uuid, f.Name)
		if f.FileInfo().IsDir() {
			os.MkdirAll(path, 0700)
		} else {
			f, err := os.OpenFile(
				path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0700)
			if err != nil {
				z.log.Errorf("Error creating file: %s", err)
			}
			defer f.Close()

			_, err = io.Copy(f, rc)
			if err != nil {
				z.log.Errorf("Error copying file contents: %s", err)
			}
		}
	}
	err = ditto(tmpUnzipDir, dst)
	if err != nil {
		z.log.Errorf("Error copying file contents: %s", err)
	}
	_ = os.RemoveAll(tmpUnzipDir)
	// move files from tmp directory to destination
	if runtimeOS := runtime.GOOS; runtimeOS == "darwin" {
		// run shell cmd
		_ = os.MkdirAll(dst, 0700)
		script := fmt.Sprintf("do shell script \"sudo xattr -d -r com.apple.quarantine '%s'\" with administrator privileges", dst)
		z.log.Infof("Move files: %s", script)
		std, err := exec.Command("osascript", "-e", script).CombinedOutput()
		if err != nil {
			z.log.Error(err)
		}
		z.log.Info(string(std))
	}

}

func (z *ZiboInstaller) GetBackups() []ZiboBackup {
	var backups []ZiboBackup
	backupDir := filepath.Join(z.Config.YazuCachePath, "backup")
	z.log.Infof("Getting backups from %s", backupDir)
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
				backupFileName := strings.ReplaceAll(backupPath, filepath.Join(z.Config.YazuCachePath, "backup"), "")
				backupFileName = strings.ReplaceAll(backupFileName, ".zip", "")
				// split by first dash
				backupFileNameSplit := strings.SplitN(backupFileName, "-", 2)
				backupVersion := strings.ReplaceAll(backupFileNameSplit[0], "/", "")
				backupDate := backupFileNameSplit[1]
				backups = append(backups, ZiboBackup{
					BackupPath: backupPath,
					Version:    backupVersion,
					Date:       backupDate,
					Size:       int(info.Size()),
				})
			}
		}
		return nil
	})
	z.log.Infof("Found %d backups", len(backups))
	return backups
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
	if installation.Path == "" {
		return
	}
	z.log.Infof("Removing %s", installation.Path)
	_ = os.RemoveAll(installation.Path)
	// if os is mac
	if runtimeOS := runtime.GOOS; runtimeOS == "darwin" {
		// run shell cmd
		script := fmt.Sprintf("do shell script \"sudo rm -rf '%s'\" with administrator privileges", installation.Path)
		std, err := exec.Command("osascript", "-e", script).CombinedOutput()
		z.log.Infof("Remov old install: %s", script)
		if err != nil {
			z.log.Error(err)
		}
		z.log.Infof(string(std))
	}
}

func (z *ZiboInstaller) DownloadZibo(fullInstall bool) (bool, string) {
	var installItem utils.Item
	var zipFilePath string
	if fullInstall {
		fullInstallItems := *z.rss.GetFullInstallItems()
		installItem = fullInstallItems[0]
	} else {
		patchedItems := *z.rss.GetPatchInstallItems()
		installItem = patchedItems[len(patchedItems)-1]
	}
	isDownloading := false

	z.log.Infof("Downloading %s, from: %s", installItem.Version, installItem.Link)
	subPath := "full"
	if !fullInstall {
		subPath = "patch/"
	}
	err := z.TorrentManager.AddTorrent(installItem.Link, subPath)
	if err != nil {
		z.log.Infof("Error downloading torrent: %s", err)
	}
	download := z.TorrentManager.Downloads[installItem.Link]
	files := download.Torrent.Files()
	file := files[0]
	zipFilePath = filepath.Join(z.TorrentManager.DownloadPath, subPath, file.Path())
	if err != nil {
		z.log.Infof("Error downloading torrent: %s", err)
	}
	isDownloading = true

	return isDownloading, zipFilePath
}

func (z *ZiboInstaller) Install(installation utils.ZiboInstallation, zipFilePath string) {
	z.unzip(zipFilePath, installation.Path, true)
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

	progress := z.TorrentManager.CheckProgress()
	return progress[link]
}

func (z *ZiboInstaller) FindInstallationDetails() utils.ZiboInstallation {
	var foundPath, version string
	res := utils.ZiboInstallation{
		Path:          foundPath,
		Version:       version,
		RemoteVersion: z.rss.GetLatestVersion(),
		BackupVersion: z.GetLastBackupVersion(),
	}
	_ = filepath.Walk(filepath.Join(z.Config.XPlanePath, "Aircraft"), func(path string, info os.FileInfo, err error) error {
		if err != nil {
			z.log.Errorf("Error walking path: %s", err)
			return err // prevent panic by handling failure accessing a path
		}
		if info.IsDir() && info.Name() == "zibomod" {
			z.log.Infof("Found zibo(mob) at: %s", path)
			foundPath = path
			return filepath.SkipDir // folder found, skip the rest of this directory
		}
		return nil
	})
	if foundPath != "" {
		foundPath = filepath.Join(foundPath, "../", "../")
		z.log.Infof("Found zibo at: %s", foundPath)
		res.Path = foundPath
		versionFilePath := filepath.Join(foundPath, "version.txt")

		data, err := os.ReadFile(versionFilePath)
		if err != nil {
			z.log.Errorf("Failed to read file: %v", err)
		}
		res.Version = string(data)
	}
	return res
}

func (z *ZiboInstaller) GetLiveries(installationDetails utils.ZiboInstallation) []InstalledLivery {
	var res []InstalledLivery
	_ = filepath.Walk(filepath.Join(installationDetails.Path, "liveries"), func(myPath string, info os.FileInfo, err error) error {
		if err != nil {
			return err // prevent panic by handling failure accessing a path
		}
		if info.IsDir() && myPath != filepath.Join(installationDetails.Path, "liveries") {
			// list png files
			_ = filepath.Walk(myPath, func(path string, myInfo os.FileInfo, err error) error {
				if myInfo.IsDir() && path != myPath {
					return filepath.SkipDir // folder found, skip the rest of this directory
				}
				if strings.LastIndex(path, "icon11.png") != -1 {
					imageBytes, err := os.ReadFile(path)
					if err != nil {
						z.log.Error(err)
					}

					// Encode the bytes to Base64
					base64Encoding := base64.StdEncoding.EncodeToString(imageBytes)
					res = append(res, InstalledLivery{
						Name: info.Name(),
						Path: myPath,
						Icon: base64Encoding,
					})
					return filepath.SkipDir // folder found, skip the rest of this directory
				}
				return nil
			})
			return filepath.SkipDir // folder found, skip the rest of this directory
		}
		return nil
	})

	return res
}

func (z *ZiboInstaller) GetAvailableLiveries() []AvailableLivery {
	var res []AvailableLivery
	//browser := rod.New().MustConnect()
	//page := browser.MustPage("https://forums.x-plane.org/index.php?/files/category/209-zibo-737/")
	//tmp := page.MustWaitLoad().MustEval("()=>Array.from(document.querySelectorAll(\"li.ipsDataItem.ipsDataItem_unread\")).map(li => li.textContent)")
	//log.Printf("%v", tmp)
	return res
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
	err = xattr.Remove(dst, "com.apple.quarantine")
	if err != nil {
		logrus.New().Warningf("failed to removing quarantine: %s", err)
	}
	return os.Chmod(dst, 0700)
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
			return os.MkdirAll(targetPath, 0700)
		} else {
			// Copy files.
			return copyFile(path, targetPath)
		}
	})
}
