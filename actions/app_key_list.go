package actions

import (
  "github.com/olekukonko/tablewriter"
  "github.com/schulterklopfer/cna/errors"
  "github.com/schulterklopfer/cna/output"
  "github.com/schulterklopfer/cna/storage"
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


