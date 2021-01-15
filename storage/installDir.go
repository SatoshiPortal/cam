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
  "fmt"
  "github.com/SatoshiPortal/cam/cyphernodeInfo"
  "github.com/SatoshiPortal/cam/dockerCompose"
  "github.com/SatoshiPortal/cam/errors"
  "github.com/SatoshiPortal/cam/globals"
  "github.com/SatoshiPortal/cam/utils"
  "github.com/SatoshiPortal/cam/version"
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

    // Check if mount point exists
    err := installedAppsIndex.MountPointHasCollision( app )
    if err != nil {
      return err
    }

    // Check if app is already installed
    apps := installedAppsIndex.Search( app.GetHash(), true )

    if len(apps) != 0 {
      return errors.APP_ALREADY_INSTALLED
    }
  }

  candidate := app.GetCandidate( version )

  if candidate == nil {
    return errors.NO_SUCH_VERSION
  }

  isRunnableErr := candidate.IsRunnableOnThisCyphernode()

  if isRunnableErr != nil {
    return isRunnableErr
  }

  err = checkAppSecurity( app, candidate )

  if err != nil {
    return err
  }

  installDirPath := utils.GetInstallDirPath()

  appHash := app.GetHash()

  err = os.MkdirAll( filepath.Join( installDirPath, appHash) , 0755 )

  if err != nil {
    return err
  }

  files :=  candidate.Files[:]
  files = append(files,globals.CANDIDATE_DESCRIPTION_FILE)

  for _, file := range files {
    sourceFilePath := filepath.Join( app.Path, globals.APP_VERSIONS_DIR, candidate.Version.Raw, file )
    targetFilePath := filepath.Join( installDirPath, appHash, file )

    if file == "docker-compose.yaml" {
      isInSwarmMode, err := cyphernodeInfo.CyphernodeIsDockerSwarm()
      if err != nil {
        return err
      }

      dockerComposeTemplate, err := dockerCompose.LoadDockerComposeTemplate( sourceFilePath, isInSwarmMode)
      if err != nil {
        return err
      }
      dockerComposeTemplate.Replacements = &map[string]string{
        "APP_UPSTREAM_HOST": app.ClientID,
        "APP_ID": app.ClientID,
        "APP_MOUNTPOINT": app.MountPoint,
      }

      // TODO: add keys and key labels to replacements
      // TODO: add mountpoint to replacements

      createTraefikLabels( dockerComposeTemplate, isInSwarmMode )
      dockerComposeTemplate.SaveAsDockerCompose( targetFilePath )
    } else {
      _, err = utils.CopyFile(sourceFilePath, targetFilePath)
      if err != nil {
        return err
      }
    }

  }
  targetFilePath := filepath.Join(installDirPath, appHash, globals.APP_DESCRIPTION_FILE)
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

func UpdateApp( app *App, version *version.Version ) error {
  installedAppsIndex, err := NewInstalledAppsIndex()
  appHash := app.GetHash()

  if err == nil {
    err = installedAppsIndex.Load()

    if err != nil {
      return err
    }
    // Check if app is already installed
    apps := installedAppsIndex.Search( appHash, true )

    if len(apps) == 0 {
      return errors.NO_SUCH_APP
    }
  }

  candidate := app.GetCandidate( version )

  if candidate == nil {
    return errors.NO_SUCH_VERSION
  }

  isRunnableErr := candidate.IsRunnableOnThisCyphernode()

  if isRunnableErr != nil {
    return isRunnableErr
  }

  err = checkAppSecurity( app, candidate )

  if err != nil {
    return err
  }

  installDirPath := utils.GetInstallDirPath()

  files :=  candidate.Files[:]
  files = append(files,globals.CANDIDATE_DESCRIPTION_FILE)

  for _, file := range files {
    sourceFilePath := filepath.Join( app.Path, globals.APP_VERSIONS_DIR, candidate.Version.Raw, file )
    targetFilePath := filepath.Join( installDirPath, appHash, file )

    if file == "docker-compose.yaml" {
      isInSwarmMode, err := cyphernodeInfo.CyphernodeIsDockerSwarm()
      if err != nil {
        return err
      }

      dockerComposeTemplate, err := dockerCompose.LoadDockerComposeTemplate( sourceFilePath, isInSwarmMode )
      if err != nil {
        return err
      }
      dockerComposeTemplate.Replacements = &map[string]string{
        "APP_UPSTREAM_HOST": app.ClientID,
        "APP_ID": app.ClientID,
        "APP_MOUNTPOINT": app.MountPoint,
      }

      createTraefikLabels( dockerComposeTemplate, isInSwarmMode )
      dockerComposeTemplate.SaveAsDockerCompose( targetFilePath )
    } else {
      _, err = utils.CopyFile(sourceFilePath, targetFilePath)
      if err != nil {
        return err
      }
    }

  }
  targetFilePath := filepath.Join(installDirPath, appHash, globals.APP_DESCRIPTION_FILE)
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

func createTraefikLabels( dockerComposeTemplate *dockerCompose.DockerComposeTemplate, isInSwarmMode bool ) {

  labels := []string{
    globals.DOCKER_COMPOSE_LABEL_TRAEFIK_ENABLE,
    globals.DOCKER_COMPOSE_LABEL_PASS_HOST_HEADER,
    globals.DOCKER_COMPOSE_LABEL_MOUNTPOINT_RULE,
    globals.DOCKER_COMPOSE_LABEL_ENTRYPOINTS,
    globals.DOCKER_COMPOSE_LABEL_ROUTER_SERVICE,
    globals.DOCKER_COMPOSE_LABEL_MW_STRIPPREXIX,
    globals.DOCKER_COMPOSE_LABEL_FORCE_SLASH,
  }

  for serviceKey, service := range *dockerComposeTemplate.Services {
    foundMainService, err := regexp.MatchString(
      "^"+fmt.Sprintf(globals.DOCKER_COMPOSE_TEMPLATE_REGEXP_TEMPLATE, "APP_UPSTREAM_HOST" )+"$" ,
      strings.Trim(serviceKey, " " ) )

    if foundMainService && err == nil {
      //var middlewares []string
      var middlewares []string //globals.DOCKER_COMPOSE_LABEL_MIDDLEWARES
      middlewareRe := regexp.MustCompile( globals.DOCKER_COMPOSE_MIDDLEWARE_PATTERN )

      for _, label := range *service.Labels {
        result := middlewareRe.FindAllStringSubmatch(label, -1)

        if result != nil && len(result) > 0 && len( result[0] ) > 1  {
          if utils.SliceIndex(len(middlewares), func( i int ) bool {
            return result[0][1] == middlewares[i]
          } ) != -1 {
            continue
          }
          middlewares = append( middlewares, result[0][1] )
        }
      }

      middlewares = append( middlewares, globals.DOCKER_COMPOSE_FORWARD_AUTH_MIDDLEWARE, globals.DOCKER_COMPOSE_STRIPPREFIX_MIDDLEWARE )

      labels = append( labels, *service.Labels... )
      labels = append( labels, globals.DOCKER_COMPOSE_LABEL_MIDDLEWARES+strings.Join(middlewares, "," ) )

      if isInSwarmMode {
        if service.Deploy == nil {
          service.Deploy = &dockerCompose.Deploy{}
        }
        service.Deploy.Labels = &labels
        service.Labels = nil
      } else {
        service.Labels = &labels
      }
      (*dockerComposeTemplate.Services)[serviceKey] = service
      break
    }
  }
}

func checkAppSecurity( app *App, candidate *AppCandidate ) error {
  if utils.SliceIndex( len(candidate.Files), func(i int) bool {
    return candidate.Files[i] == "docker-compose.yaml"
  } ) > -1 {

    isInSwarmMode, err := cyphernodeInfo.CyphernodeIsDockerSwarm()
    if err != nil {
      return err
    }

    dockerComposeTemplate, err := dockerCompose.LoadDockerComposeTemplate(
      filepath.Join(app.Path, globals.APP_VERSIONS_DIR, candidate.Version.Raw, "docker-compose.yaml" ),
      isInSwarmMode )

    if err != nil {
      return err
    }

    err = dockerComposeTemplate.CheckVolumes( app.TrustZone )

    if err != nil {
      return err
    }

    err = dockerComposeTemplate.CheckServiceNames()

    if err != nil {
      return err
    }

    err = dockerComposeTemplate.CheckNetworks( app.TrustZone, app.ClientID )

    if err != nil {
      return err
    }
  }
  return nil
}
