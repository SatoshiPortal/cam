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
  "github.com/SatoshiPortal/cam/globals"
  "github.com/SatoshiPortal/cam/output"
  "github.com/SatoshiPortal/cam/storage"
  "github.com/SatoshiPortal/cam/utils"
  "github.com/SatoshiPortal/cam/version"
  "github.com/urfave/cli"
  "regexp"
  "sort"
  "strings"
)

func App_install(c *cli.Context) error {
  if len(c.Args()) == 0 {
    return errors.APP_INSTALL_NO_APP_ID
  }

  repoIndex, err := storage.NewRepoIndex()

  if err != nil {
    return err
  }

  err = repoIndex.Load()

  if err != nil {
    return err
  }

  isTrusted_appToInstall_Version := strings.Trim( c.Args().Get(0), " \n")

  trustZone := globals.DefaultTrustZone

  appToInstall_Version := strings.Split( isTrusted_appToInstall_Version, ":" )

  var appToInstall string

  if len(appToInstall_Version) > 1 {
    trustZone = strings.ToLower(appToInstall_Version[0])
    if utils.SliceIndex( len(globals.ValidTrustZones), func(i int) bool {
      return globals.ValidTrustZones[i] == trustZone
    } ) == -1 {
      output.Noticef( "Unknown trust zone \"%s\". Valid values are: %s\n", trustZone, strings.Join( globals.ValidTrustZones, ", ") )
      return nil
    }
    appToInstall = appToInstall_Version[1]
  } else {
    appToInstall = appToInstall_Version[0]
  }

  appIDVersion := strings.Split( appToInstall, "@" )
  v := version.NewVersion( "latest" )
  appHandle := appIDVersion[0]
  if len(appIDVersion) > 1 {
    v.Parse( appIDVersion[1] )
  }

  apps := repoIndex.Search( appHandle, false )
  sort.Slice(apps, func(i, j int) bool {
    return apps[i].Label < apps[j].Label
  })

  if len(apps) == 0 {
    return errors.NO_SUCH_APP
  } else if len(apps) == 1 {
    if v.Raw == "latest" {
      v = version.NewVersion( apps[0].Latest )
    }

    mountPoint := apps[0].DefaultMountPoint
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
      mountPoint = apps[0].DefaultMountPoint
      currentMountPointAttempt = mountPoint
      mountPoint = sanitizeMountPoint(mountPoint)
      output.Noticef( "unable to use \"%s\" as mount point. Switching to app hash as mount point\n", currentMountPointAttempt)
      mountPoint = apps[0].GetHash()
    }

    apps[0].MountPoint = mountPoint
    apps[0].Secret = utils.RandomString(32, base32.StdEncoding.EncodeToString )
    apps[0].ClientID = apps[0].GetHash()
    apps[0].TrustZone = trustZone
    err := storage.InstallApp( apps[0], v )

    if err != nil {
      return err
    }

    output.Noticef( "installed \"%s\" at version %s from %s at /%s\n", apps[0].Label, v.Raw, apps[0].Source.String(), apps[0].MountPoint )

  } else {
    output.Notice( "Multiple apps with that label exist. Please install one using the hash." )
    for i:=0; i<len(apps); i++ {
      output.Noticef( "%s - %s (%24s)\n", apps[i].Label, apps[i].Source.String(), apps[i].GetHash() )
    }
  }
  return nil
}

func sanitizeMountPoint( mountPoint string ) string {
  for strings.HasPrefix( mountPoint, "/" ) {
    mountPoint = mountPoint[1:]
  }
  mountPointReplacementsRegexp := regexp.MustCompile(`[^a-zA-Z0-9-_]+`)
  return mountPointReplacementsRegexp.ReplaceAllLiteralString( mountPoint, "_" )
}