/*
 * MIT License
 *
 * Copyright (c) 2020 schulterklopfer/__escapee__
 *
 * Permission is hereby granted, free of charge, to any person obtaining a copy
 * of this software and associated documentation files (the "Software"), to deal
 * in the Software without restriction, including without limitation the rights
 * to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
 * copies of the Software, and to permit persons to whom the Software is
 * furnished to do so, subject to the following conditions:
 *
 * The above copyright notice and this permission notice shall be included in all
 * copies or substantial portions of the Software.
 *
 * THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
 * IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
 * FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
 * AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
 * LIABILIT * Y, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
 * OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
 * SOFTWARE.
 */

package storage

import (
  "encoding/json"
  errors2 "errors"
  "github.com/SatoshiPortal/cam/errors"
  "github.com/SatoshiPortal/cam/globals"
  "github.com/SatoshiPortal/cam/utils"
  "io/ioutil"
  "os"
  "path/filepath"
)

type InstalledAppsIndex struct {
  AppList
}

func NewInstalledAppsIndex() (*InstalledAppsIndex, error) {
  if !utils.InstalledAppsIndexFileExists() {
    return &InstalledAppsIndex{}, errors.INSTALLED_APPS_INDEX_DOES_NOT_EXIST
  }
  return &InstalledAppsIndex{ AppList{ Apps: []*App{}, Labels: map[string][]int{} } }, nil
}

func (installedAppsIndex *InstalledAppsIndex) Load() error {
  if !utils.InstalledAppsIndexFileExists() {
    return errors.INSTALLED_APPS_INDEX_DOES_NOT_EXIST
  }

  installedAppsIndexJsonBytes, err := ioutil.ReadFile( utils.GetInstalledAppsIndexFilePath() )
  if err != nil {
    return err
  }

  err = json.Unmarshal(installedAppsIndexJsonBytes, &installedAppsIndex)
  if err != nil {
    return err
  }
  installedAppsIndex.BuildAppHashes()
  installedAppsIndex.BuildLabels()
  return nil
}

func (installedAppsIndex *InstalledAppsIndex) Build() error {

  if exists, err := utils.InstallDirExists(); err != nil || !exists {
    return errors.INSTALL_DIR_DOES_NOT_EXIST
  }
  // get contents of app installation dir:

  installDir := utils.GetInstallDirPath()

  d, err := os.Open(installDir)

  if err != nil {
    if d != nil {
      _ = d.Close()
    }
    return errors2.New( "Could not process install dir: "+ err.Error() )
  }

  files, err := d.Readdir(-1)
  if err != nil {
    if d != nil {
      _ = d.Close()
    }
    return errors2.New( "Could not process install dir: "+ err.Error() )
  }

  installedAppsIndex.Clear()

  for _, file := range files {
    if !file.IsDir() {
      continue
    }
    appDescriptionPath := filepath.Join(installDir,file.Name(),globals.APP_DESCRIPTION_FILE)

    appDescriptionJsonBytes, err := ioutil.ReadFile( appDescriptionPath )
    if err != nil {
      continue
    }
    var app App

    err = json.Unmarshal( appDescriptionJsonBytes, &app )
    if err != nil {
      continue
    }

    app.Path = filepath.Join(installDir,file.Name())

    candidateDescriptionPath := filepath.Join(installDir, file.Name(), globals.CANDIDATE_DESCRIPTION_FILE)

    candidateDescriptionJsonBytes, err := ioutil.ReadFile( candidateDescriptionPath )
    if err != nil {
      continue
    }

    var candidate AppCandidate
    err = json.Unmarshal( candidateDescriptionJsonBytes, &candidate )

    if err != nil {
      continue
    }

    app.Candidates = []*AppCandidate{ &candidate }
    app.BuildHash()
    err = installedAppsIndex.AddApp( &app )

    if err != nil {
      continue
    }
  }

  installedAppsIndex.BuildLabels()

  return installedAppsIndex.Save()

}

func (installedAppsIndex *InstalledAppsIndex) Save() error {
  installedAppsIndexJsonBytes, err := json.MarshalIndent( installedAppsIndex, "", "  " )
  err = ioutil.WriteFile(utils.GetInstalledAppsIndexFilePath(), installedAppsIndexJsonBytes, 0644)

  if err != nil {
    return err
  }
  return nil
}

func (installedAppsIndex *InstalledAppsIndex) MountPointHasCollision( app *App ) error {
  // check if mountPoint is the same or below any other installed app's mount point
  for _,a := range installedAppsIndex.Apps {
    if a.MountPoint == app.MountPoint {
      if a.ClientID == app.ClientID {
        return errors.APP_ALREADY_INSTALLED
      } else {
        return errors.APP_MOUNTPOINT_BLOCKED
      }
    }
  }
  return nil
}