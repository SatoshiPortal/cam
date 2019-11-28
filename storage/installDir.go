package storage

import (
  "encoding/base32"
  "encoding/json"
  "github.com/schulterklopfer/cam/dockerCompose"
  "github.com/schulterklopfer/cam/errors"
  "github.com/schulterklopfer/cam/globals"
  "github.com/schulterklopfer/cam/utils"
  "github.com/schulterklopfer/cam/version"
  "io/ioutil"
  "os"
  "path/filepath"
)

func InitInstallDir() error {
  installDir := utils.GetInstallDirPath()

  err := os.MkdirAll(installDir, 0755)
  if err != nil {
    return err
  }

  installedAppsIndex, err := NewInstalledAppsIndex()

  if err != nil {
    err = installedAppsIndex.Build()

    if err != nil {
      return err
    }
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

  isRunnable, _ := AppCandidateIsRunnableOnCyphernode( candidate )

  if !isRunnable {
    return errors.APP_VERSION_IS_NOT_COMPATIBLE
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

    if file == "docker-compose.yaml" {
      dockerComposeTemplate, err := dockerCompose.LoadDockerComposeTemplate( sourceFilePath )
      if err != nil {
        return err
      }
      dockerComposeTemplate.Replacements = &map[string]string{
        "APP_UPSTREAM_HOST": app.ClientID,
        "APP_ID": app.ClientID,
      }
      dockerComposeTemplate.SaveAsDockerCompose( targetFilePath )
    } else {
      _, err = utils.CopyFile(sourceFilePath, targetFilePath)
      if err != nil {
        return err
      }
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

func AddKeyToApp( app *App, key *Key ) error {
  if app == nil || key == nil {
    return nil
  }

  if utils.SliceIndex( len(app.Keys), func(i int) bool {
    return app.Keys[i].Label == key.Label
  } ) != -1 {
    return nil
  }

  targetFilePath := filepath.Join( utils.GetInstallDirPath(), app.GetHash(), globals.APP_DESCRIPTION_FILE)

  app.Keys = append( app.Keys, key )
  appDescriptionJsonBytes, err := json.MarshalIndent( app, "", "  " )

  if err != nil {
    return err
  }

  err = ioutil.WriteFile(targetFilePath, appDescriptionJsonBytes, 0644)

  if err != nil {
    return err
  }

  return nil
}

func RemoveKeyFromApp( app *App, key *Key ) error {
  if app == nil || key == nil {
    return nil
  }

  keyIndex := utils.SliceIndex( len(app.Keys), func(i int) bool {
    return app.Keys[i].Label == key.Label
  } )

  if keyIndex == -1 {
    return nil
  }

  targetFilePath := filepath.Join( utils.GetInstallDirPath(), app.GetHash(), globals.APP_DESCRIPTION_FILE)

  app.Keys = append(app.Keys[:keyIndex], app.Keys[keyIndex+1:]...)
  appDescriptionJsonBytes, err := json.MarshalIndent( app, "", "  " )

  if err != nil {
    return err
  }

  err = ioutil.WriteFile(targetFilePath, appDescriptionJsonBytes, 0644)

  if err != nil {
    return err
  }

  return nil
}

func checkAppSecurity( app *App, candidate *AppCandidate ) error {
  if utils.SliceIndex( len(candidate.Files), func(i int) bool {
    return candidate.Files[i] == "docker-compose.yaml"
  } ) > -1 {

    dockerComposeTemplate, err := dockerCompose.LoadDockerComposeTemplate(
      filepath.Join(app.Path, globals.APP_VERSIONS_DIR, candidate.Version.Raw, "docker-compose.yaml" ) )

    if err != nil {
      return err
    }

    err = dockerComposeTemplate.CheckVolumes()

    if err != nil {
      return err
    }

    err = dockerComposeTemplate.CheckServiceNames()

    if err != nil {
      return err
    }
  }
  return nil
}
