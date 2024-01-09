package main

import (
	"changeme/installer"
	"changeme/utils"
	"context"
	"github.com/wailsapp/wails/v2/pkg/runtime"
	"log"
	"os"
	"path/filepath"
)

// App struct
type App struct {
	ctx    context.Context
	config *installer.Config
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{
		config: installer.NewConfig(utils.RealHomeDirGetter{}),
	}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
}

func (a *App) IsXPlanePathConfigured() bool {
	return a.config.CheckXPlanePath(a.config.XPlanePath)
}
func (a *App) CheckXPlanePath(dirPath string) bool {
	return a.config.CheckXPlanePath(dirPath)
}
func (a *App) GetConfig() installer.Config {
	return *a.config
}
func (a *App) OpenDirDialog() string {
	res, _ := runtime.OpenDirectoryDialog(a.ctx, runtime.OpenDialogOptions{})
	return res
}

func (a *App) FindZiboInstallationDetails() installer.ZiboInstallation {
	// find the zibo folder in the X-Plane directory
	var foundPath, version string
	_ = filepath.Walk(filepath.Join(a.config.XPlanePath, "aircraft"), func(path string, info os.FileInfo, err error) error {
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
	rss := utils.NewRss("https://skymatixva.com/tfiles/feed.xml")
	return installer.ZiboInstallation{
		Path:          foundPath,
		Version:       version,
		RemoteVersion: rss.GetLatestVersion(),
		BackupVersion: a.config.GetLastBackupVersion(),
	}
}

func (a *App) BackupZiboInstallation(installation installer.ZiboInstallation) bool {
	return a.config.Backup(installation)
}

func (a *App) RestoreZiboInstallation(installation installer.ZiboInstallation) bool {
	if installation.Version == "" {
		installation.Path = filepath.Join(a.config.XPlanePath, "Aircraft", "B737-800X")
	}
	return a.config.Restore(installation)
}

func (a *App) FreshInstallZibo(installation installer.ZiboInstallation) {
	if installation.Version == "" {
		installation.Path = filepath.Join(a.config.XPlanePath, "Aircraft", "B737-800X")
	} else {
		a.config.RemoveOldInstalls(installation)
	}
	a.config.Install(installation)
}

func (a *App) DownloadZibo(fullInstall bool) bool {
	return a.config.DownloadZibo(fullInstall)
}

func (a *App) UpdateZibo(installation installer.ZiboInstallation) {
	a.config.Update(installation)
}

func (a *App) GetDownloadDetails(update bool) float64 {
	return a.config.GetDownloadProgress(update)
}
