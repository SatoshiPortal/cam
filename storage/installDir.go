package storage

import (
  "encoding/base32"
  "encoding/json"
  "github.com/schulterklopfer/cna/dockerCompose"
  "github.com/schulterklopfer/cna/errors"
  "github.com/schulterklopfer/cna/globals"
  "github.com/schulterklopfer/cna/utils"
  "github.com/schulterklopfer/cna/version"
  "gopkg.in/yaml.v2"
  "io/ioutil"
  "os"
  "path/filepath"
  "regexp"
  "strings"
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

func InstallApp( app *App, version *version.Version ) error {
  installedAppsIndex, err := NewInstalledAppsIndex()

  if err == nil {
    err = installedAppsIndex.Load()

    if err != nil {
      return err
    }

    apps := installedAppsIndex.Search( app.GetHash(), true )

    if len(apps) != 0 {
      return errors.APP_ALREADY_INSTALLED
    }
  }

  candidate := app.GetCandidate( version )

  if candidate == nil {
    return errors.NO_SUCH_VERSION
  }

  err = checkAppSecurity( app, candidate )

  if err != nil {
    return err
  }

  installDirPath := utils.GetInstallDirPath()

  clientID := app.GetHash()

  err = os.MkdirAll( filepath.Join( installDirPath, clientID ) , 0755 )

  if err != nil {
    return err
  }

  app.ClientSecret = utils.RandomString(32, base32.StdEncoding.EncodeToString )
  app.ClientID = clientID

  files :=  candidate.Files[:]
  files = append(files,globals.CANDIDATE_DESCRIPTION_FILE)

  for _, file := range files {
    sourceFilePath := filepath.Join( app.Path, globals.APP_VERSIONS_DIR, candidate.Version.Raw, file )
    targetFilePath := filepath.Join( installDirPath, clientID, file )
    _, err = utils.CopyFile(sourceFilePath, targetFilePath)
    if err != nil {
      return err
    }
  }
  targetFilePath := filepath.Join(installDirPath, clientID, globals.APP_DESCRIPTION_FILE)
  appDescriptionJsonBytes, err := json.MarshalIndent( app, "", "  " )

  err = ioutil.WriteFile(targetFilePath, appDescriptionJsonBytes, 0644)

  if err != nil {
    return err
  }

  err = installedAppsIndex.Build()

  if err != nil {
    return err
  }

  return nil
}

func UninstallApp( app *App ) error {
  installedAppsIndex, err := NewInstalledAppsIndex()

  if err == nil {
    err = installedAppsIndex.Load()

    if err != nil {
      return err
    }

    apps := installedAppsIndex.Search( app.GetHash(), true )

    if len(apps) != 1 {
      return errors.APP_NOT_INSTALLED
    }
  }

  installDirPath := utils.GetInstallDirPath()

  clientID := app.GetHash()

  err = os.RemoveAll( filepath.Join(installDirPath,clientID) )

  if err != nil {
    return err
  }

  err = installedAppsIndex.Build()

  if err != nil {
    return err
  }

  return nil
}


func checkAppSecurity( app *App, candidate *AppCandidate ) error {
  if utils.SliceIndex( len(candidate.Files), func(i int) bool {
    return candidate.Files[i] == "docker-compose.yaml"
  } ) > -1 {
    err := checkDockerCompose( filepath.Join(app.Path, globals.APP_VERSIONS_DIR, candidate.Version.Raw, "docker-compose.yaml" ) )
    if err != nil {
      return err
    }
  }
  return nil
}

func checkDockerCompose( path string ) error {
  dockerComposeBytes, err := ioutil.ReadFile(path)
  if err != nil {
    return err
  }
  var dockerCompose dockerCompose.Config

  err = yaml.Unmarshal(dockerComposeBytes, &dockerCompose)
  if err != nil {
    return err
  }

  for _, service := range dockerCompose.Services {
    for _, volume := range service.Volumes {
      arr := strings.Split( volume, ":" )
      hostDirectory := strings.Trim( arr[0], " \n" )
      if utils.SliceIndex( len(globals.DockerVolumeWhitelist), func(i int) bool {
        pattern := globals.DockerVolumeWhitelist[i]
        match, err := regexp.MatchString(pattern, hostDirectory)
        return match && err == nil
      } ) == -1 {
        return errors.VOLUME_NOT_IN_WHITELIST
      }

      if utils.SliceIndex( len(globals.DockerVolumeElementBlacklist), func(i int) bool {
        return strings.Contains( hostDirectory, globals.DockerVolumeElementBlacklist[i] )
      } ) != -1 {
        return errors.VOLUME_HAS_ILLEGAL_ELEMENTS
      }
    }
  }

  return nil
}

