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
