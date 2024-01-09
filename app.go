package main

import (
	"changeme/installer"
	"changeme/utils"
	"context"
	"github.com/wailsapp/wails/v2/pkg/runtime"
	"path/filepath"
)

// App struct
type App struct {
	ctx  context.Context
	zibo *installer.ZiboInstaller
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{
		zibo: installer.NewZibo(utils.RealHomeDirGetter{}),
	}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
}

func (a *App) IsXPlanePathConfigured() bool {
	config := utils.GetConfig(utils.RealHomeDirGetter{})
	return config.CheckXPlanePath(a.zibo.Config.XPlanePath)
}
func (a *App) CheckXPlanePath(dirPath string) bool {
	config := utils.GetConfig(utils.RealHomeDirGetter{})
	return config.CheckXPlanePath(dirPath)
}
func (a *App) GetConfig() utils.Config {
	return *a.zibo.Config
}
func (a *App) OpenDirDialog() string {
	res, _ := runtime.OpenDirectoryDialog(a.ctx, runtime.OpenDialogOptions{})
	return res
}

func (a *App) FindZiboInstallationDetails() utils.ZiboInstallation {
	// find the zibo folder in the X-Plane directory
	return a.zibo.FindInstallationDetails()
}

func (a *App) BackupZiboInstallation(installation utils.ZiboInstallation) bool {
	return a.zibo.Backup(installation)
}

func (a *App) RestoreZiboInstallation(installation utils.ZiboInstallation) bool {
	if installation.Version == "" {
		installation.Path = filepath.Join(a.zibo.Config.XPlanePath, "Aircraft", "B737-800X")
	}
	return a.zibo.Restore(installation)
}

func (a *App) InstallZibo(installation utils.ZiboInstallation) {
	if installation.Version == "" {
		installation.Path = filepath.Join(a.zibo.Config.XPlanePath, "Aircraft", "B737-800X")
	} else {
		a.zibo.RemoveOldInstalls(installation)
	}
	a.zibo.Install(installation)
}

func (a *App) DownloadZibo(fullInstall bool) bool {
	return a.zibo.DownloadZibo(fullInstall)
}

func (a *App) UpdateZibo(installation utils.ZiboInstallation) {
	a.zibo.Update(installation)
}

func (a *App) GetDownloadDetails(update bool) float64 {
	return a.zibo.GetDownloadProgress(update)
}
