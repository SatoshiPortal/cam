package commands

import (
  "github.com/schulterklopfer/cam/actions"
  "github.com/urfave/cli"
)

func InitKeyCommands( app *cli.App ) {
  app.Commands = append( app.Commands, cli.Command{
    Name:    "key",
    Aliases: []string{"k"},
    Usage:   "key commands",
    Subcommands: []cli.Command{
      {
        Name:    "list",
        Aliases: []string{"l"},
        Usage:   "lists keys",
        Action: actions.ActionWrapper(actions.Key_list, false, true),
      },
    },
  },
  )
}
