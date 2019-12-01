package commands

import (
  "github.com/SatoshiPortal/cam/actions"
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
        Usage:   "list installed apps",
        Action: actions.ActionWrapper(actions.App_list),
      },
      {
        Name:    "install",
        Aliases: []string{"i"},
        Usage:   "installs an app",
        Action: actions.ActionWrapper(actions.App_install),
      },
      {
        Name:    "delete",
        Aliases: []string{"d"},
        Usage:   "deletes an app",
        Action: actions.ActionWrapper(actions.App_delete),
      },
      {
        Name:    "search",
        Aliases: []string{"s"},
        Usage:   "search for an app in all the sources",
        Action: actions.ActionWrapper(actions.App_search),
      },
      {
        Name:    "key",
        Aliases: []string{"k"},
        Usage:   "handle cypherapp keys",
        Subcommands: []cli.Command{
          {
            Name:    "list",
            Aliases: []string{"l"},
            Usage:   "list keys for an app",
            Action: actions.ActionWrapper(actions.App_keyList),
          },
          {
            Name:    "add",
            Aliases: []string{"a"},
            Usage:   "add key to app",
            Action: actions.ActionWrapper(actions.App_keyAdd),
          },
          {
            Name:    "delete",
            Aliases: []string{"d"},
            Usage:   "delete key from app",
            Action: actions.ActionWrapper(actions.App_keyDelete),
          },
        },
      },
    },
  },
  )
}
