package actions

import (
  "github.com/schulterklopfer/cam/errors"
  "github.com/schulterklopfer/cam/storage"
  "github.com/schulterklopfer/cam/utils"
  "github.com/urfave/cli"
)

func ActionWrapper( action func(c *cli.Context) error, boolParams ...bool ) func(c *cli.Context) error {
  needsDataDir := true
  needsKeysFile := false
  needsLock := true
  if len(boolParams) > 0 {
    needsDataDir = boolParams[0]
  }
  if len(boolParams) > 1 {
    needsKeysFile = boolParams[1]
  }
  if len(boolParams) > 2 {
    needsLock = boolParams[2]
  }
  return func(c *cli.Context) error {
    if needsKeysFile && !utils.KeysFileExists() {
      return errors.NO_KEYS_FILE
    }

    if needsDataDir && !utils.ValidDataDirExists() {
      return errors.DATADIR_IS_INVALID
    }

    // check for write access
    if needsLock {
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
