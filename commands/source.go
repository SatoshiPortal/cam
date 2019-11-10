package commands

import (
  "github.com/schulterklopfer/cna/actions"
  "github.com/urfave/cli"
)

func InitSourceCommands( app *cli.App ) {
  app.Commands = append( app.Commands, cli.Command{
    Name:    "source",
    Aliases: []string{"s"},
    Usage:   "source commands",
    Subcommands: []cli.Command{
      {
        Name:    "list",
        Aliases: []string{"l"},
        Usage:   "lists sources",
        Action: actions.ActionWrapper(actions.Source_list),
      },
      {
        Name:    "add",
        Aliases: []string{"a"},
        Usage:   "add source",
        Action: actions.ActionWrapper(actions.Source_add, true),
      },
      {
        Name:    "del",
        Aliases: []string{"d"},
        Usage:   "delete source",
        Action: actions.ActionWrapper(actions.Source_delete,true),
      },
    },
  },
  )
}
