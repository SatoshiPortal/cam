package actions

import (
  "github.com/olekukonko/tablewriter"
  "github.com/schulterklopfer/cna/storage"
  "github.com/urfave/cli"
  "os"
  "sort"
)

func App_list(c *cli.Context) error {

  installedAppsIndex, err := storage.NewInstalledAppsIndex()

  if err != nil {
    return err
  }

  err = installedAppsIndex.Load()

  if err != nil {
    return err
  }

  apps := installedAppsIndex.Apps[:]
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