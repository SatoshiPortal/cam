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
  "github.com/SatoshiPortal/cam/errors"
  "github.com/SatoshiPortal/cam/storage"
  "github.com/SatoshiPortal/cam/utils"
  "github.com/urfave/cli"
)

// TODO: add command to copy keys.properties to .cam

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
