package actions

import (
  "github.com/olekukonko/tablewriter"
  "github.com/schulterklopfer/cam/storage"
  "github.com/urfave/cli"
  "os"
  "sort"
  "strings"
)

func Key_list(c *cli.Context) error {
  keyList := storage.NewKeyList()
  err := keyList.Load()
  if err != nil {
    return err
  }

  keys := keyList.Keys[:]
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

  return nil
}

