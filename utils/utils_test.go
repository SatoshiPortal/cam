package utils_test

import (
  "github.com/schulterklopfer/cam/globals"
  "github.com/schulterklopfer/cam/utils"
  "os"
  "path/filepath"
  "testing"
)

func TestDataDirExists(t *testing.T) {

  _ = os.RemoveAll( globals.DATA_DIR )
  defer os.RemoveAll( globals.DATA_DIR )

  dataDirExists, _ := utils.DataDirExists()

  if dataDirExists {
    t.Error( "expecting: dataDirExists==false" )
  }

  _,_ = os.Create( globals.DATA_DIR )

  dataDirExists, err := utils.DataDirExists()

  if !dataDirExists || err == nil {
    t.Error( "expecting: dataDirExists==true with error" )
  }

  _ = os.RemoveAll( globals.DATA_DIR )
  _ = os.Mkdir( globals.DATA_DIR, 0777 )

  dataDirExists, err = utils.DataDirExists()

  if !dataDirExists {
    t.Error( "expecting: dataDirExists==true" )
  }

}

func TestIsLocked(t *testing.T) {
  _ = os.RemoveAll( globals.DATA_DIR )
  _ = os.Mkdir( globals.DATA_DIR, 0777 )
  _,_ = os.Create( filepath.Join( globals.DATA_DIR, globals.LOCK_FILE ) )

  defer os.RemoveAll( globals.DATA_DIR )

  isLocked, err := utils.IsLocked()

  if !isLocked || err != nil {
    t.Error( "expecting: isLocked==true with no error" )
  }

}

func TestLock(t *testing.T) {
  _ = os.RemoveAll( globals.DATA_DIR )
  _ = os.Mkdir( globals.DATA_DIR, 0777 )

  defer os.RemoveAll( globals.DATA_DIR )

  err := utils.Lock()

  if err != nil  {
    t.Error( "expecting: no error" )
  }

  _, err = os.Stat( filepath.Join(globals.DATA_DIR, globals.LOCK_FILE) )

  if err != nil  {
    t.Error( "expecting: no error" )
  }
}

func TestUnlock(t *testing.T) {
  _ = os.RemoveAll( globals.DATA_DIR )
  _ = os.Mkdir( globals.DATA_DIR, 0777 )
  _,_ = os.Create( filepath.Join( globals.DATA_DIR, globals.LOCK_FILE ) )

  defer os.RemoveAll( globals.DATA_DIR )

  err := utils.Unlock()

  if err != nil  {
    t.Error( "expecting: no error" )
  }

  _, err = os.Stat( filepath.Join(globals.DATA_DIR, globals.LOCK_FILE) )

  if err == nil  {
    t.Error( "expecting: error" )
  }
}

func TestLockFileExists(t *testing.T) {
  _ = os.RemoveAll( globals.DATA_DIR )
  _ = os.Mkdir( globals.DATA_DIR, 0777 )
  _,_ = os.Create( filepath.Join( globals.DATA_DIR, globals.LOCK_FILE ) )

  defer os.RemoveAll( globals.DATA_DIR )

  lockFileExists, err := utils.LockFileExists()

  if !lockFileExists || err != nil  {
    t.Error( "expecting: lockFileExists == true and no error" )
  }
}

func TestStateFileExists(t *testing.T) {
  _ = os.RemoveAll( globals.DATA_DIR )
  _ = os.Mkdir( globals.DATA_DIR, 0777 )
  _,_ = os.Create( filepath.Join( globals.DATA_DIR, globals.STATE_FILE ) )

  defer os.RemoveAll( globals.DATA_DIR )

  stateFileExists, err := utils.StateFileExists()

  if !stateFileExists || err != nil  {
    t.Error( "expecting: stateFileExists == true and no error" )
  }
}

func TestInitDataDir(t *testing.T) {
  _ = os.RemoveAll( globals.DATA_DIR )
  defer os.RemoveAll( globals.DATA_DIR )

  err := utils.InitDataDir()

  if err != nil  {
    t.Error( "expecting: no error" )
  }

  //_, err = os.Stat( filepath.Join(globals.DATA_DIR, globals.LOCK_FILE) )

  //if err != nil  {
  //  t.Error( "expecting: no error" )
  //}

}

func TestInitStateFile(t *testing.T) {
  _ = os.RemoveAll( globals.DATA_DIR )
  _ = os.Mkdir( globals.DATA_DIR, 0777 )
  _,_ = os.Create( filepath.Join( globals.DATA_DIR, globals.STATE_FILE ) )

  defer os.RemoveAll( globals.DATA_DIR )

  _, err := utils.InitStateFile()

  if err != nil {
    t.Error(  "expecting: no error" )
  }

  _, err = os.Stat( filepath.Join(globals.DATA_DIR, globals.STATE_FILE) )

  if err != nil  {
    t.Error( "expecting: no error" )
  }

}
