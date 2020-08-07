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
  "github.com/SatoshiPortal/cam/errors"
  "github.com/SatoshiPortal/cam/output"
  "github.com/SatoshiPortal/cam/storage"
  "github.com/urfave/cli"
  "sort"
  "strings"
)

func App_delete(c *cli.Context) error {
  if len(c.Args()) == 0 {
    return errors.APP_INSTALL_NO_APP_ID
  }

  installedAppsIndex, err := storage.NewInstalledAppsIndex()

  if err != nil {
    return err
  }

  err = installedAppsIndex.Load()

  if err != nil {
    return err
  }

  appHandle := strings.Trim( c.Args().Get(0), " \n")

  apps := installedAppsIndex.Search( appHandle, false )
  sort.Slice(apps, func(i, j int) bool {
    return apps[i].Label < apps[j].Label
  })

  if len(apps) == 0 {
    return errors.NO_SUCH_APP
  } else if len(apps) == 1 {
    println( "deleting \""+apps[0].Label+"\"" )
    err := storage.UninstallApp( apps[0] )
    if err != nil {
      return err
    }
  } else {
    println( "Multiple apps matching that label are installed. Please try deleting using the hash." )
    for i:=0; i<len(apps); i++ {
      output.Noticef( "%s (%24s)\n", apps[i].Label, apps[i].GetHash() )
    }
  }
  return nil
}