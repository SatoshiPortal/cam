package storage

import (
  "github.com/schulterklopfer/cna/utils"
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
