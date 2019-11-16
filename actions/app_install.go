package actions

import (
  "github.com/schulterklopfer/cna/errors"
  "github.com/schulterklopfer/cna/output"
  "github.com/schulterklopfer/cna/storage"
  "github.com/schulterklopfer/cna/version"
  "github.com/urfave/cli"
  "sort"
  "strings"
)

func App_install(c *cli.Context) error {
  // TODO: upgrade without removal
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

  appToInstall := strings.Trim( c.Args().Get(0), " \n")

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
    println( "installing \""+apps[0].Label+"\" at version "+v.Raw+" from "+apps[0].Source.String() )
    err := storage.InstallApp( apps[0], v )
    if err != nil {
      return err
    }
  } else {
    println( "Multiple apps with that label exist. Please install one using the hash." )
    for i:=0; i<len(apps); i++ {
      output.Noticef( "%s - %s (%24s)\n", apps[i].Label, apps[i].Source.String(), apps[i].GetHash() )
    }
  }
  return nil
}