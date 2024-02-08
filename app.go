package main

import (
	"context"
	"encoding/json"
	"github.com/sirupsen/logrus"
	"github.com/wailsapp/wails/v2/pkg/runtime"
	"github/xairline/yazu/installer"
	"github/xairline/yazu/utils"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	goruntime "runtime"
)

// App struct
type App struct {
	ctx  context.Context
	zibo *installer.ZiboInstaller
	Log  *logrus.Logger
}

type DownloadInfo struct {
	IsDownloading bool   `json:"isDownloading"`
	Path          string `json:"path"`
}

// NewApp creates a new App application struct
func NewApp() *App {
	log := logrus.New()
	zibo := installer.NewZibo(utils.RealHomeDirGetter{}, true, log)
	homeDir, _ := utils.RealHomeDirGetter{}.UserHomeDir()
	file, err := os.OpenFile(
		filepath.Join(homeDir, ".yazu", "yazu.log"),
		os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Error("Failed to log to file, using default stderr")
	}
	log.SetOutput(io.MultiWriter(file, os.Stdout))

	return &App{
		zibo: zibo,
		Log:  log,
	}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
}

func (a *App) IsXPlanePathConfigured() bool {
	config := utils.GetConfig(utils.RealHomeDirGetter{}, true, a.Log)
	return config.CheckXPlanePath(a.zibo.Config.XPlanePath, []string{})
}
func (a *App) CheckXPlanePath(dirPath string, cachePath []string) bool {
	config := utils.GetConfig(utils.RealHomeDirGetter{}, true, a.Log)
	return config.CheckXPlanePath(dirPath, cachePath)
}
func (a *App) GetConfig() utils.Config {
	return *a.zibo.Config
}
func (a *App) OpenDirDialog() string {
	res, _ := runtime.OpenDirectoryDialog(a.ctx, runtime.OpenDialogOptions{})
	return res
}

func (a *App) FindZiboInstallationDetails() []utils.ZiboInstallation {
	// find the zibo folder in the X-Plane directory
	return a.zibo.FindInstallationDetails()
}

func (a *App) BackupZiboInstallation(installation utils.ZiboInstallation) bool {
	_, err := a.zibo.Backup(installation)
	return err != nil
}

func (a *App) RestoreZiboInstallation(installation utils.ZiboInstallation, backupPath string) bool {
	if installation.Version == "" {
		installation.Path = filepath.Join(a.zibo.Config.XPlanePath, "Aircraft", "B737-800X")
	}
	err := a.zibo.Restore(installation, backupPath)
	return err != nil
}

func (a *App) InstallZibo(installation utils.ZiboInstallation, zipPath string) {
	if installation.Version == "" {
		installation.Path = filepath.Join(a.zibo.Config.XPlanePath, "Aircraft", "B737-800X")
	} else {
		a.zibo.RemoveOldInstalls(installation)
	}
	a.zibo.Install(installation, zipPath)
}

func (a *App) DownloadZibo(fullInstall bool) DownloadInfo {
	isDownloading, zipFilePath := a.zibo.DownloadZibo(fullInstall)
	a.Log.Printf("isDownloading: %v, zipFilePath: %v", isDownloading, zipFilePath)
	res := DownloadInfo{
		IsDownloading: isDownloading,
		Path:          zipFilePath,
	}
	return res
}

func (a *App) UpdateZibo(installation utils.ZiboInstallation, zipPath string) {
	a.zibo.Update(installation, zipPath)
}

func (a *App) GetDownloadDetails(update bool) float64 {
	return a.zibo.GetDownloadProgress(update)
}

func (a *App) GetBackups() []installer.ZiboBackup {
	return a.zibo.GetBackups()
}

func (a *App) GetCachedFiles() []utils.CachedFile {
	return a.zibo.TorrentManager.GetCachedFiles()
}

func (a *App) GetLiveries(installation utils.ZiboInstallation) []installer.InstalledLivery {
	return a.zibo.GetLiveries(installation)
}

func (a *App) GetAvailableLiveries() []installer.AvailableLivery {
	// return a.zibo.GetAvailableLiveries()
	return []installer.AvailableLivery{}
}

func (a *App) GetOs() string {
	return goruntime.GOOS
}

func (a *App) DeleteFiles(files []string) string {
	for _, file := range files {
		err := os.Remove(file)
		if err != nil {
			return err.Error()
		}
	}
	return ""
}

func (a *App) GetVersion() string {
	return AppVersion
}

func (a *App) GetLatestVersion() string {
	type GitHubRelease struct {
		TagName string `json:"tag_name"` // The name of the tag for this release
	}
	url := "https://api.github.com/repos/xairline/yazu/releases/latest"
	resp, err := http.Get(url)
	if err != nil {
		return "unknown"
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "unknown"
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "unknown"
	}

	var release GitHubRelease
	err = json.Unmarshal(body, &release)
	if err != nil {
		return "unknown"
	}

	return release.TagName
}
