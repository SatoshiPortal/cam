package storage_test

import (
  "github.com/SatoshiPortal/cam/globals"
  "github.com/SatoshiPortal/cam/storage"
  "io/ioutil"
  "os"
  "path/filepath"
  "testing"
  "time"
)

func TestLoadStateFile(t *testing.T) {
  jsonStateString := `{ "version": "test", "lastUpdate": "`+time.Now().Format(time.RFC3339)+`" }`
  defer os.RemoveAll( globals.DATA_DIR )
  _ = os.RemoveAll( globals.DATA_DIR )
  _ = os.Mkdir( globals.DATA_DIR, 0777 )
  _ = ioutil.WriteFile( filepath.Join( globals.DATA_DIR, globals.STATE_FILE ) , []byte(jsonStateString), 0644)

  _, err := storage.LoadStateFile( filepath.Join( globals.DATA_DIR, globals.STATE_FILE ) )

  if err != nil {
    t.Error( "expected: no error")
  }
}

func TestStateFile_Save(t *testing.T) {
  defer os.RemoveAll( globals.DATA_DIR )
  _ = os.RemoveAll( globals.DATA_DIR )
  _ = os.Mkdir( globals.DATA_DIR, 0777 )

  stateFile := storage.NewStateFile( filepath.Join( globals.DATA_DIR, globals.STATE_FILE ) )
  stateFile.Save()
  _, err := os.Stat( filepath.Join(globals.DATA_DIR, globals.STATE_FILE) )

  if err != nil  {
    t.Error( "expecting: no error" )
  }

}
