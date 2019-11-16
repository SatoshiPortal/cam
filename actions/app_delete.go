package actions

import (
  "github.com/schulterklopfer/cam/errors"
  "github.com/schulterklopfer/cam/output"
  "github.com/schulterklopfer/cam/storage"
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