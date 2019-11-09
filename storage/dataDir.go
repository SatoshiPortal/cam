package storage

import (
  "github.com/schulterklopfer/cna/utils"
  "os"
)

func InitDataDir() error {
  dataDir, err := utils.GetDataDirPath()
  if err != nil {
    return err
  }
  err = os.MkdirAll( dataDir, 0755)
  if err != nil {
    return err
  }
  repoDir, err := utils.GetRepoDirPath()
  if err != nil {
    return err
  }
  err = os.MkdirAll( repoDir, 0755)
  if err != nil {
    return err
  }

  // check if state file exists:
  stateFileExists, _ := utils.StateFileExists()
  if !stateFileExists {
    stateFilePath, err := utils.GetStateFilePath()
    if err != nil {
      return err
    }
    stateFile := NewStateFile( stateFilePath )
    err = stateFile.Save()
    if err != nil {
      return err
    }
  }

  // check if source file exists:
  sourceFileExists, _ := utils.SourceFileExists()
  if !sourceFileExists {
    sourceFilePath, err := utils.GetSourceFilePath()
    if err != nil {
      return err
    }
    sourceFile := NewSourceList( sourceFilePath )
    err = sourceFile.Save()
    if err != nil {
      return err
    }
  }

  return nil
}

func Lock() error {
  lockFilePath, err := utils.GetLockFilePath()
  if err != nil {
    return err
  }
  _, err = os.Create( lockFilePath )
  if err != nil {
    return err
  }
  return nil
}

func Unlock() error {
  lockFilePath, err := utils.GetLockFilePath()
  if err != nil {
    return err
  }
  err = os.RemoveAll( lockFilePath )
  if err != nil {
    return err
  }
  return nil
}

func IsLocked() (bool, error) {
  return utils.LockFileExists()
}

