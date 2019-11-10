package storage

import (
  "encoding/json"
  "github.com/schulterklopfer/cna/errors"
  "github.com/schulterklopfer/cna/globals"
  "github.com/schulterklopfer/cna/output"
  "github.com/schulterklopfer/cna/utils"
  "github.com/schulterklopfer/cna/version"
  "io/ioutil"
  "os"
  "path/filepath"
)

type RepoIndex struct {
  Apps *AppList `json:"apps"`
}

func NewRepoIndex() (*RepoIndex, error) {
  if !utils.RepoIndexFileExists() {
    return &RepoIndex{}, errors.REPO_INDEX_DOES_NOT_EXIST
  }
  return &RepoIndex{},nil
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
  return nil
}

func (repoIndex *RepoIndex) Build() error {

  sourceList, err := LoadSourceFile( utils.GetSourceFilePath() )

  if err != nil {
    return err
  }

  appList := new( AppList )

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
        candidate.Version = version.NewVersion(versionsDFile.Name())
        candidates = append( candidates, &candidate )
      }

      app.Candidates = candidates
      app.BuildHash()
      err = appList.AddApp( &app )

      if err != nil {
        continue
      }
    }
  }

  repoIndexJsonBytes, err := json.MarshalIndent( appList, "", "  " )
  err = ioutil.WriteFile(utils.GetRepoIndexFilePath(), repoIndexJsonBytes, 0644)
  if err != nil {
    return err
  }

  return nil

}
