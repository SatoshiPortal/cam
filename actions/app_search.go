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
  "github.com/SatoshiPortal/cam/storage"
  "github.com/urfave/cli"
  "os"
  "sort"
  "strings"
)

func App_search(c *cli.Context) error {

  if len(c.Args()) == 0 {
    return errors.APP_SEARCH_NO_SEARCH_TERM
  }

  repoIndex, err := storage.NewRepoIndex()

  if err != nil {
    return err
  }

  err = repoIndex.Load()

  if err != nil {
    return err
  }

  searchString := strings.Trim( c.Args().Get(0), " \n")

  apps := repoIndex.Search( searchString, false )
  sort.Slice(apps, func(i, j int) bool {
    return apps[i].Label < apps[j].Label
  })

  table := tablewriter.NewWriter(os.Stdout)
  table.SetHeader([]string{"Source", "Label", "Name", "Hash"})
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

  for i:=0; i<len(apps); i++ {
    table.Append( []string{ apps[i].Source.String(),apps[i].Label,apps[i].Name,apps[i].GetHash()} )
  }

  table.Render() // Send output

  return nil

}
