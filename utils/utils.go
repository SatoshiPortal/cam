package utils

import (
  "encoding/base64"
  "github.com/schulterklopfer/cna/errors"
  "github.com/schulterklopfer/cna/globals"
  "golang.org/x/crypto/ripemd160"
  "os"
  "path/filepath"
  "strings"
)

func fileExists( path string ) bool {
  if _, err := os.Stat( path ); err != nil {
    return false
  }
  return true
}

func dirExists( path string ) (bool, error) {
  fileInfo, err := os.Stat( path )
  if  err != nil  {
    return false,err
  }
  if !fileInfo.IsDir() {
    return true, errors.DIR_IS_NOT_A_DIRECTORY
  }
  return true, nil
}

func fileExistsInDataDir( file string ) bool {
  return fileExists( filepath.Join( GetDataDirPath(), file ) )
}

func dirExistsInDataDir( d string ) (bool, error) {
  return dirExists( getPath( d ) )
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

func GetInstallDirPath() string {
  installDir := os.Getenv( globals.INSTALL_DIR_ENV_KEY )
  if installDir  == "" {
    cwd, err := os.Getwd()
    if err != nil {
      panic( err )
    }
    installDir = filepath.Join( cwd, "install" )
  }
  return installDir
}

func DataDirExists() (bool, error) {
  return dirExistsInDataDir( "" )
}

func InstallDirExists() (bool, error) {
  return dirExists( GetInstallDirPath() )
}

func RepoDirExists() (bool, error) {
  return dirExistsInDataDir( globals.REPO_DIR )
}

func RepoExists( repoDir string ) (bool,error) {
  return dirExistsInDataDir( filepath.Join(globals.REPO_DIR,repoDir) )
}

func StateFileExists() bool {
  return fileExistsInDataDir( globals.STATE_FILE )
}

func InstalledAppsFileExists() bool {
  return fileExists( GetInstallDirPath() )
}

func SourceFileExists() bool {
  return fileExistsInDataDir( globals.SOURCE_FILE )
}

func LockFileExists() bool {
  return fileExistsInDataDir( globals.LOCK_FILE )
}

func RepoIndexFileExists() bool {
  return fileExistsInDataDir( globals.REPO_INDEX_FILE )
}

func GetInstalledAppsIndexFilePath() string {
  return filepath.Join( GetInstallDirPath(), globals.INSTALLED_APPS_FILE )
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
  hasher := ripemd160.New()
  hasher.Write(*bytes)
  return strings.Trim(base64.URLEncoding.EncodeToString(hasher.Sum(nil)), "=" )
}

