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

package storage

import (
  "github.com/SatoshiPortal/cam/utils"
  "os"
)

func InitDataDir() error {
  dataDir := utils.GetDataDirPath()

  err := os.MkdirAll( dataDir, 0755)
  if err != nil {
    return err
  }
  repoDir := utils.GetRepoDirPath()

  err = os.MkdirAll( repoDir, 0755)
  if err != nil {
    return err
  }

  // check if state file exists:
  if !utils.StateFileExists() {
    stateFile := NewStateFile( utils.GetStateFilePath() )
    err = stateFile.Save()
    if err != nil {
      return err
    }
  }

  // check if source file exists:
  if !utils.SourceFileExists() {
    sourceFile := NewSourceList( utils.GetSourceFilePath() )
    err = sourceFile.Save()
    if err != nil {
      return err
    }
  }

  return nil
}

func Lock() error {
  // Wait for lockedfile Mutex decision of golang dev
  _, err := os.Create( utils.GetLockFilePath() )
  if err != nil {
    return err
  }
  return nil
}

func Unlock() error {
  // Wait for lockedfile Mutex decision of golang dev
  err := os.RemoveAll( utils.GetLockFilePath() )
  if err != nil {
    return err
  }
  return nil
}

func IsLocked() bool {
  // Wait for lockedfile Mutex decision of golang dev
  return utils.LockFileExists()
}
