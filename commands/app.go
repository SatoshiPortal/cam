package commands

import (
  "github.com/schulterklopfer/cna/actions"
  "github.com/urfave/cli"
)

func InitAppCommands( app *cli.App ) {
  app.Commands = append( app.Commands, cli.Command{
    Name:    "app",
    Aliases: []string{"a"},
    Usage:   "app (list|update)",
    Subcommands: []cli.Command{
      {
        Name:    "list",
        Aliases: []string{"l"},
        Usage:   "lists installed apps",
        Action: actions.ActionWrapper(actions.App_list),
      },
      {
        Name:    "update",
        Aliases: []string{"l"},
        Usage:   "upates app index using all sources",
        Action: actions.ActionWrapper(actions.App_update,true),
      },
    },
  },
  )
}
