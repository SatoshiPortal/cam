package commands

import (
  "github.com/schulterklopfer/cna/actions"
  "github.com/urfave/cli"
)

func InitAppCommands( app *cli.App ) {
  app.Commands = append( app.Commands, cli.Command{
    Name:    "app",
    Aliases: []string{"a"},
    Usage:   "app commands",
    Subcommands: []cli.Command{
      {
        Name:    "list",
        Aliases: []string{"l"},
        Usage:   "lists available apps",
        Action: actions.ActionWrapper(actions.App_list),
      },
      {
        Name:    "search",
        Aliases: []string{"s"},
        Usage:   "search for an app in all the sources",
        Action: actions.ActionWrapper(actions.App_search, true),
      },
    },
  },
  )
}
