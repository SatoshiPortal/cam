package utils

import (
  "github.com/schulterklopfer/cna/errors"
  "github.com/schulterklopfer/cna/globals"
  "os"
  "path/filepath"
)

func fileExists( file string ) (bool,error) {
  dataDirExists, err := DataDirExists()

  if err != nil {
    return false,err
  }

  if !dataDirExists {
    return false, errors.DATADIR_DOES_NOT_EXIST
  }

  dataDir, err := GetDataDirPath()
  if err != nil {
    return false, err
  }

  _, err = os.Stat(filepath.Join( dataDir, file ))
  if  err != nil  {
    return false,err
  }

  return true,nil
}

func dirExists( d string ) (bool, error) {
  dir, err := getPath( d )
  if err != nil {
    return false, err
  }

  fileInfo, err := os.Stat(dir)
  if  err != nil  {
    return false,err
  }

  if !fileInfo.IsDir() {
    return true, errors.DATADIR_IS_NOT_A_DIRECTORY
  }
  return true, nil
}

func getPath( file string ) (string,error) {
  dataDir, err := GetDataDirPath()
  if err != nil {
    return "", err
  }
  return filepath.Join( dataDir, file ),nil
}

func GetDataDirPath() (string, error) {
  cwd, err := os.Getwd()
  if err != nil {
    return "",err
  }
  return filepath.Join( cwd, globals.DATA_DIR ), nil
}

func DataDirExists() (bool, error) {
  return dirExists( "" )
}

func RepoDirExists() (bool, error) {
  return dirExists( globals.REPO_DIR )
}

func RepoExists( repoDir string ) (bool,error) {
  return dirExists( filepath.Join(globals.REPO_DIR,repoDir) )
}

func StateFileExists() (bool, error) {
  return fileExists( globals.STATE_FILE )
}

func SourceFileExists() (bool, error) {
  return fileExists( globals.SOURCE_FILE )
}

func LockFileExists() (bool, error) {
  return fileExists( globals.LOCK_FILE )
}

func GetStateFilePath() (string, error) {
  return getPath(globals.STATE_FILE)
}

func GetLockFilePath() (string, error) {
  return getPath(globals.LOCK_FILE)
}

func GetSourceFilePath() (string, error) {
  return getPath(globals.SOURCE_FILE)
}

func GetRepoDirPath() (string, error) {
  return getPath(globals.REPO_DIR)
}

func GetRepoDirPathFor( repo string ) (string, error) {
  repoDir, err := GetRepoDirPath()
  if err != nil {
    return "",err
  }
  return filepath.Join(repoDir,repo), nil
}

func ValidDataDirExists() bool {
  dataDirExists, err := DataDirExists()
  if dataDirExists && err != nil {
    return false
  }
  if !dataDirExists {
    return false
  }

  if stateFileExists, _ := StateFileExists(); !stateFileExists {
    return false
  }

  if sourceFileExists, _ := StateFileExists(); !sourceFileExists {
    return false
  }

  if repoDirExists, _ := RepoDirExists(); !repoDirExists {
    return false
  }

  return true

}
