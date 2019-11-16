package actions

import (
  "github.com/schulterklopfer/cam/errors"
  "github.com/schulterklopfer/cam/output"
  "github.com/schulterklopfer/cam/storage"
  "github.com/urfave/cli"
  "strings"
)

func App_keyDelete(c *cli.Context) error {
  if len(c.Args()) < 2 {
    return errors.APP_SEARCH_NO_SEARCH_TERM
  }

  installedAppsIndex, err := storage.NewInstalledAppsIndex()

  if err != nil {
    return err
  }

  err = installedAppsIndex.Load()

  if err != nil {
    return err
  }

  searchString := strings.Trim( c.Args().Get(0), " \n")
  keyLabel := strings.Trim( c.Args().Get(1), " \n")

  apps := installedAppsIndex.Search( searchString, false )
  if len(apps) == 0 {
    return errors.NO_SUCH_APP
  } else if len(apps) == 1 {
    println( "deleting key from \""+apps[0].Label+"\"" )

    keyList := storage.NewKeyList()
    err := keyList.Load()
    if err != nil {
      return err
    }

    key := keyList.GetKey( keyLabel )
    if key == nil {
      return errors.NO_SUCH_KEY
    }

    err = storage.RemoveKeyFromApp( apps[0], key )

    if err != nil {
      return err
    }

  } else {
    println( "Multiple apps with that label exist. Please use hash instead." )
    for i:=0; i<len(apps); i++ {
      output.Noticef( "%s - %s (%24s)\n", apps[i].Label, apps[i].Source.String(), apps[i].GetHash() )
    }
  }

  err = installedAppsIndex.Build()
  if err != nil {
    return err
  }

  return nil
}

