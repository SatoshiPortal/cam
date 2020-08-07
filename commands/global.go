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
