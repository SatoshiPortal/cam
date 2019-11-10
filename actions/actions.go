package actions

import (
  "github.com/schulterklopfer/cna/errors"
  "github.com/schulterklopfer/cna/storage"
  "github.com/schulterklopfer/cna/utils"
  "github.com/urfave/cli"
)

func ActionWrapper( action func(c *cli.Context) error, boolParams ...bool ) func(c *cli.Context) error {
  needsDataDir := true
  needsWriteAccess := false
  if len(boolParams) > 0 {
    needsWriteAccess = boolParams[0]
  }
  if len(boolParams) > 1 {
    needsDataDir = boolParams[1]
  }
  return func(c *cli.Context) error {
    // check if we have a repository
    if needsDataDir && !utils.ValidDataDirExists() {
      return errors.DATADIR_IS_INVALID
    }

    // check for write access
    if needsWriteAccess {
      if storage.IsLocked() {
        return errors.DATADIR_IS_LOCKED
      }
      err := storage.Lock()
      defer storage.Unlock()
      if err != nil {
        return err
      }
    }
    return action( c )
  }
}
