package actions

import (
  "github.com/SatoshiPortal/cam/errors"
  "github.com/SatoshiPortal/cam/storage"
  "github.com/SatoshiPortal/cam/utils"
  "github.com/urfave/cli"
  "strings"
)

func Source_delete(c *cli.Context) error {
  sourceList, err := storage.LoadSourceFile( utils.GetSourceFilePath() )

  if err != nil {
    return err
  }

  if len(c.Args()) == 0 {
    return errors.SOURCE_DELETE_NO_SOURCE
  }

  sourceString := strings.Trim( c.Args().Get(0), " \n")
  err = sourceList.RemoveSource( sourceString )
  if err != nil {
    return err
  }

  err = sourceList.Save()
  if err != nil {
    return err
  }

  return nil
}
