package utils

import (
  "crypto/rand"
  "encoding/base64"
  "fmt"
  "github.com/schulterklopfer/cam/errors"
  "github.com/schulterklopfer/cam/globals"
  "golang.org/x/crypto/ripemd160"
  "io"
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

func GetKeysFilePath() string {
  fromEnv := os.Getenv( globals.KEYS_FILE_ENV_KEY )
  if fromEnv == "" {
    return filepath.Join( GetDataDirPath(), "keys.properties" )
  }
  return fromEnv
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

func InstalledAppsIndexFileExists() bool {
  return fileExists( filepath.Join(GetInstallDirPath(), globals.INSTALLED_APPS_FILE ) )
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

func KeysFileExists() bool {
  return fileExists( GetKeysFilePath() )
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

func RandomString(length int, encodeToString func([]byte) string ) string {
  randomBytes := make([]byte, length)
  if _, err := io.ReadFull(rand.Reader, randomBytes); err != nil {
    return ""
  }
  return strings.TrimRight( encodeToString( randomBytes), "=" )
}

func SliceIndex(limit int, predicate func(i int) bool) int {
  for i := 0; i < limit; i++ {
    if predicate(i) {
      return i
    }
  }
  return -1
}

func CopyFile(src string, dst string) (int64, error) {
  sourceFileStat, err := os.Stat(src)
  if err != nil {
    return 0, err
  }

  if !sourceFileStat.Mode().IsRegular() {
    return 0, fmt.Errorf("%s is not a regular file", src)
  }

  source, err := os.Open(src)
  if err != nil {
    return 0, err
  }
  defer source.Close()

  destination, err := os.Create(dst)
  if err != nil {
    return 0, err
  }
  defer destination.Close()
  nBytes, err := io.Copy(destination, source)
  return nBytes, err
}

func AddIndexToLabel( dict *map[string][]int, label string, index int ) {
  if _, ok := (*dict)[label]; !ok {
    (*dict)[label] = make( []int, 0 )
  }

  if SliceIndex( len((*dict)[label]), func(i int) bool {
    return (*dict)[label][i] == index
  } ) == -1 {
    (*dict)[label] = append((*dict)[label], index)
  }

}