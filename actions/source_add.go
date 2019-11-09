package actions

import (
  "github.com/schulterklopfer/cna/errors"
  "github.com/schulterklopfer/cna/storage"
  "github.com/schulterklopfer/cna/utils"
  "github.com/urfave/cli"
  "strings"
)

func Source_add(c *cli.Context) error {
  sourceFilePath, err := utils.GetSourceFilePath()
  if err != nil {
    return err
  }
  sourceList, err := storage.LoadSourceFile( sourceFilePath )

  if err != nil {
    return err
  }

  if len(c.Args()) == 0 {
    return errors.SOURCE_ADD_NO_SOURCE
  }

  sourceString := strings.Trim( c.Args().Get(0), " \n")
  err = sourceList.AddSource( sourceString )
  if err != nil {
    return err
  }

  err = sourceList.Save()
  if err != nil {
    return err
  }

  return nil
}
