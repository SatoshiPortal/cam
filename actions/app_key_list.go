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

package actions

import (
  "github.com/olekukonko/tablewriter"
  "github.com/SatoshiPortal/cam/errors"
  "github.com/SatoshiPortal/cam/output"
  "github.com/SatoshiPortal/cam/storage"
  "github.com/urfave/cli"
  "os"
  "sort"
  "strings"
)

func App_keyList(c *cli.Context) error {

  if len(c.Args()) == 0 {
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

  apps := installedAppsIndex.Search( searchString, false )
  if len(apps) == 0 {
    return errors.NO_SUCH_APP
  } else if len(apps) == 1 {
    println( "keys for \""+apps[0].Label+"\"" )

    keys := apps[0].Keys[:]
    sort.Slice(keys, func(i, j int) bool {
      return len(keys[i].Groups) < len(keys[j].Groups)
    })

    table := tablewriter.NewWriter(os.Stdout)
    table.SetHeader([]string{"Label", "Data", "Groups"})
    table.SetAutoFormatHeaders(true)
    table.SetHeaderAlignment(tablewriter.ALIGN_LEFT)
    table.SetAlignment(tablewriter.ALIGN_LEFT)
    table.SetCenterSeparator("")
    table.SetColumnSeparator("")
    table.SetRowSeparator("")
    table.SetHeaderLine(false)
    table.SetBorder(false)
    table.SetTablePadding("  ")
    table.SetNoWhiteSpace(true)

    for i:=0; i<len(keys); i++ {
      table.Append( []string{ keys[i].Label,keys[i].Data,strings.Join(keys[i].Groups, ", " )} )
    }

    table.Render() // Send output

  } else {
    println( "Multiple apps with that label exist. Please use hash instead." )
    for i:=0; i<len(apps); i++ {
      output.Noticef( "%s - %s (%24s)\n", apps[i].Label, apps[i].Source.String(), apps[i].GetHash() )
    }
  }

  return nil
}


