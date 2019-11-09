package actions

import (
  "github.com/schulterklopfer/cna/output"
  "github.com/schulterklopfer/cna/storage"
  "github.com/schulterklopfer/cna/utils"
  "github.com/urfave/cli"
)

func Source_list(c *cli.Context) error {
  sourceFilePath, err := utils.GetSourceFilePath()
  if err != nil {
    return err
  }
  sourceList, err := storage.LoadSourceFile( sourceFilePath )

  if err != nil {
    return err
  }

  for i:=0; i<len(sourceList.Sources); i++ {
    output.Noticef( "* %s\n", sourceList.Sources[i].String() )
  }

  return nil
}
