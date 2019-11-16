package commands

import (
  "github.com/schulterklopfer/cam/actions"
  "github.com/urfave/cli"
)

func InitGlobalCommands( app *cli.App ) {
  app.Commands = append( app.Commands, []cli.Command{
    {
      Name:    "init",
      Aliases: []string{"i"},
      Usage:   "inits the storage folder in the current directory",
      Action:  actions.ActionWrapper( actions.Global_init, false, false, false ),
    },
    {
      Name:    "update",
      Aliases: []string{"u"},
      Usage:   "updates local app repositories from the sources",
      Action:  actions.ActionWrapper( actions.Global_update ),
    },
  }...
  )
}
