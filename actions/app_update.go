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

package actions

import (
  "encoding/base32"
  "github.com/SatoshiPortal/cam/errors"
  "github.com/SatoshiPortal/cam/output"
  "github.com/SatoshiPortal/cam/storage"
  "github.com/SatoshiPortal/cam/utils"
  "github.com/SatoshiPortal/cam/version"
  "github.com/urfave/cli"
  "strings"
)

func App_update(c *cli.Context) error {

  if len(c.Args()) == 0 {
    return errors.APP_INSTALL_NO_APP_ID
  }

  repoIndex := storage.NewRepoIndex()

  err := repoIndex.Load()

  if err != nil {
    return err
  }

  appToUpdate := strings.Trim( c.Args().Get(0), " \n")

  appIDVersion := strings.Split(appToUpdate, "@" )
  v := version.NewVersion( "latest" )
  appHandle := appIDVersion[0]
  if len(appIDVersion) > 1 {
    v.Parse( appIDVersion[1] )
  }

  installedAppsIndex, err := storage.NewInstalledAppsIndex()

  if err != nil {
    return err
  }

  err = installedAppsIndex.Load()

  if err != nil {
    return err
  }

  repoApps := repoIndex.Search( appHandle, false )

  if len(repoApps) == 0 {
    return errors.NO_SUCH_APP
  } else {

    var apps []*storage.App

    for _,a := range repoApps {
      installedApps := installedAppsIndex.Search( a.GetHash(), true )
      apps = append( apps, installedApps... )
    }

    if len(apps) == 1 {
      appIndex := repoIndex.AppIndex( apps[0] )
      app := repoIndex.Apps[appIndex]
      println( "updating \""+app.Label+"\"" )
      //    err := storage.UninstallApp( apps[0] )
      //    if err != nil {
      //      return err
      //    }
      if v.Raw == "latest" {
        v = version.NewVersion( app.Latest )
      }

      mountPoint := app.DefaultMountPoint
      currentMountPointAttempt := ""

      if len(c.Args()) > 1 {
        mountPoint = strings.Trim( c.Args().Get(1), " \n")
        currentMountPointAttempt = mountPoint
        mountPoint = sanitizeMountPoint(mountPoint)
        if mountPoint == "" {
          output.Noticef( "unable to use \"%s\" as mount point. Switching to default mount point\n", currentMountPointAttempt)
        }
      }

      if mountPoint == "" {
        mountPoint = app.DefaultMountPoint
        currentMountPointAttempt = mountPoint
        mountPoint = sanitizeMountPoint(mountPoint)
        output.Noticef( "unable to use \"%s\" as mount point. Switching to app hash as mount point\n", currentMountPointAttempt)
        mountPoint = apps[0].GetHash()
      }

      app.MountPoint = mountPoint
      app.Secret = utils.RandomString(32, base32.StdEncoding.EncodeToString )
      app.ClientID = app.GetHash()
      app.TrustZone = apps[0].TrustZone
      err := storage.UpdateApp( app, v )

      if err != nil {
        return err
      }

      output.Noticef( "updated \"%s\" to version %s from %s at /%s\n", app.Label, v.Raw, app.Source.String(), app.MountPoint )

    } else {
      println( "Multiple apps matching that label are installed. Please try updating using the hash." )
      for i:=0; i<len(apps); i++ {
        output.Noticef( "%s (%24s)\n", apps[i].Label, apps[i].GetHash() )
      }
    }
  }

  return nil

}