package utils

import (
  "encoding/base64"
  "github.com/schulterklopfer/cna/errors"
  "github.com/schulterklopfer/cna/globals"
  "golang.org/x/crypto/ripemd160"
  "hash"
  "os"
  "path/filepath"
  "strings"
)

var hasher hash.Hash

func fileExists( file string ) bool {
  dataDirExists, err := DataDirExists()

  if err != nil {
    return false
  }

  if !dataDirExists {
    return false
  }

  dataDir := GetDataDirPath()

  _, err = os.Stat(filepath.Join( dataDir, file ))
  if  err != nil  {
    return false
  }

  return true
}

func dirExists( d string ) (bool, error) {
  dir := getPath( d )

  fileInfo, err := os.Stat(dir)
  if  err != nil  {
    return false,err
  }

  if !fileInfo.IsDir() {
    return true, errors.DATADIR_IS_NOT_A_DIRECTORY
  }
  return true, nil
}

func getPath( file string ) string {
  return filepath.Join( GetDataDirPath(), file )
}

func GetDataDirPath() string {
  cwd, err := os.Getwd()
  if err != nil {
    panic( err )
  }
  return filepath.Join( cwd, globals.DATA_DIR )
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

func StateFileExists() bool {
  return fileExists( globals.STATE_FILE )
}

func SourceFileExists() bool {
  return fileExists( globals.SOURCE_FILE )
}

func LockFileExists() bool {
  return fileExists( globals.LOCK_FILE )
}

func RepoIndexFileExists() bool {
  return fileExists( globals.REPO_INDEX_FILE )
}

func GetStateFilePath() string {
  return getPath(globals.STATE_FILE)
}

func GetLockFilePath() string {
  return getPath(globals.LOCK_FILE)
}

func GetSourceFilePath() string {
  return getPath(globals.SOURCE_FILE)
}

func GetRepoIndexFilePath() string {
  return getPath( globals.REPO_INDEX_FILE )
}

func GetRepoDirPath() string {
  return getPath(globals.REPO_DIR)
}

func GetRepoDirPathFor( repo string ) string {
  return filepath.Join( GetRepoDirPath(),repo)
}

func ValidDataDirExists() bool {
  dataDirExists, err := DataDirExists()
  if dataDirExists && err != nil {
    return false
  }
  if !dataDirExists {
    return false
  }

  if !StateFileExists() {
    return false
  }

  if !SourceFileExists() {
    return false
  }

  if repoDirExists, _ := RepoDirExists(); !repoDirExists {
    return false
  }

  return true

}

func BuildHash( bytes *[]byte ) string {
  if hasher == nil {
    hasher = ripemd160.New()
  }
  hasher.Write(*bytes)
  return strings.Trim(base64.URLEncoding.EncodeToString(hasher.Sum(nil)), "=" )
}

