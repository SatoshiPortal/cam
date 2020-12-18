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
  "github.com/SatoshiPortal/cam/errors"
  "github.com/SatoshiPortal/cam/globals"
  "github.com/SatoshiPortal/cam/output"
  "github.com/SatoshiPortal/cam/utils"
  "io/ioutil"
  "os"
  "path/filepath"
)

type RepoIndex struct {
  AppList
}

func NewRepoIndex() *RepoIndex {
  return &RepoIndex{}
}


func (repoIndex *RepoIndex) Load() error {
  if !utils.RepoIndexFileExists() {
    return errors.REPO_INDEX_DOES_NOT_EXIST
  }

  repoIndexJsonBytes, err := ioutil.ReadFile( utils.GetRepoIndexFilePath() )
  if err != nil {
    return err
  }

  err = json.Unmarshal( repoIndexJsonBytes, &repoIndex )
  if err != nil {
    return err
  }
  repoIndex.BuildAppHashes()
  repoIndex.BuildLabels()
  return nil
}

func (repoIndex *RepoIndex) Build() error {

  sourceList, err := LoadSourceFile( utils.GetSourceFilePath() )

  if err != nil {
    return err
  }
  repoIndex.Clear()

  for i:=0; i<len(sourceList.Sources); i++ {

    absolutePath := sourceList.Sources[i].GetAbsolutePath()

    d, err := os.Open(absolutePath)

    if err != nil {
      if d != nil {
        _ = d.Close()
      }
      output.Warningf( "Could not process source %s: %s\n", sourceList.Sources[i].String(), err.Error() )
      continue
    }

    files, err := d.Readdir(-1)
    if err != nil {
      if d != nil {
        _ = d.Close()
      }
      output.Warningf( "Could not process source %s: %s\n", sourceList.Sources[i].String(), err.Error() )
      continue
    }

    _ = d.Close()

    for _, file := range files {
      if !file.IsDir() {
        continue
      }

      appDescriptionPath := filepath.Join(absolutePath,file.Name(),globals.APP_DESCRIPTION_FILE)

      appDescriptionJsonBytes, err := ioutil.ReadFile( appDescriptionPath )
      if err != nil {
        continue
      }

      var app App

      err = json.Unmarshal( appDescriptionJsonBytes, &app )
      if err != nil {
        continue
      }

      app.Path = filepath.Join(absolutePath,file.Name())
      app.Source = sourceList.Sources[i]

      appVersionsPath := filepath.Join(app.Path, globals.APP_VERSIONS_DIR)

      versionsD, err := os.Open(appVersionsPath)

      if err != nil {
        if versionsD != nil {
          _ = versionsD.Close()
        }
        continue
      }

      versionsDFiles, err := versionsD.Readdir(-1)
      if err != nil {
        if versionsD != nil {
          _ = versionsD.Close()
        }
        continue
      }

      _ = versionsD.Close()

      candidates := make( []*AppCandidate,0 )
      for _, versionsDFile := range versionsDFiles {
        if !versionsDFile.IsDir() {
          continue
        }

        candidateDescriptionPath := filepath.Join(appVersionsPath,versionsDFile.Name(), globals.CANDIDATE_DESCRIPTION_FILE)

        candidateDescriptionJsonBytes, err := ioutil.ReadFile( candidateDescriptionPath )
        if err != nil {
          continue
        }

        var candidate AppCandidate
        err = json.Unmarshal( candidateDescriptionJsonBytes, &candidate )
        if err != nil {
          continue
        }
        if candidate.Version.Raw != versionsDFile.Name() {
          output.Warningf( "App candidate version mismatch. Ignoring %s@%s from %s\n", app.Label, candidate.Version.Raw, app.Source.String()  )
        } else {
          candidates = append( candidates, &candidate )
        }
      }

      app.Candidates = candidates
      app.BuildHash()
      err = repoIndex.AddApp( &app )

      if err != nil {
        continue
      }
    }
  }
  repoIndex.BuildLabels()

  return repoIndex.Save()

}

func (repoIndex *RepoIndex) Save() error {
  repoIndexJsonBytes, err := json.MarshalIndent( repoIndex, "", "  " )
  err = ioutil.WriteFile(utils.GetRepoIndexFilePath(), repoIndexJsonBytes, 0644)
  if err != nil {
    return err
  }
  return nil
}
