package storage

import (
  "github.com/schulterklopfer/cna/errors"
  "github.com/schulterklopfer/cna/utils"
  "github.com/schulterklopfer/cna/version"
  "os"
)

func InitInstallDir() error {
  installDir := utils.GetInstallDirPath()

  err := os.MkdirAll(installDir, 0755)
  if err != nil {
    return err
  }

  installedAppsIndex, err := NewInstalledAppsIndex()

  if err != nil {
    return err
  }

  err = installedAppsIndex.Build()

  if err != nil {
    return err
  }

  return nil
}

func InstallApp( app *App, v *version.Version ) error {
  installedAppsIndex, err := NewInstalledAppsIndex()

  if err != nil {
    return err
  }

  err = installedAppsIndex.Load()

  if err != nil {
    return err
  }

  apps := installedAppsIndex.Search( app.GetHash(), true )

  if len(apps) != 0 {
    return errors.APP_ALREADY_INSTALLED
  }

  return nil
}

