package actions

import (
  "github.com/olekukonko/tablewriter"
  "github.com/schulterklopfer/cam/errors"
  "github.com/schulterklopfer/cam/storage"
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
