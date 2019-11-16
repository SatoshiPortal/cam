package storage_test

import (
  "github.com/schulterklopfer/cam/globals"
  "github.com/schulterklopfer/cam/storage"
  "io/ioutil"
  "os"
  "path/filepath"
  "testing"
)

func TestLoadSourceFile(t *testing.T) {
  sourcesString := `git://github.com/SatoshiPortal/cypherapps.git
file:///home/jash/test`
  defer os.RemoveAll( globals.DATA_DIR )
  _ = os.RemoveAll( globals.DATA_DIR )
  _ = os.Mkdir( globals.DATA_DIR, 0777 )
  _ = ioutil.WriteFile( filepath.Join( globals.DATA_DIR, globals.SOURCE_FILE ) , []byte(sourcesString), 0644)

  _, err := storage.LoadSourceFile( filepath.Join( globals.DATA_DIR, globals.SOURCE_FILE ) )

  if err != nil {
    t.Error( "expected: no error")
  }
}

func TestSourceFile_Save(t *testing.T) {
  defer os.RemoveAll( globals.DATA_DIR )
  _ = os.RemoveAll( globals.DATA_DIR )
  _ = os.Mkdir( globals.DATA_DIR, 0777 )

  sourceFile := storage.NewSourceList( filepath.Join( globals.DATA_DIR, globals.SOURCE_FILE ) )
  sourceFile.Save()
  _, err := os.Stat( filepath.Join(globals.DATA_DIR, globals.SOURCE_FILE) )

  if err != nil  {
    t.Error( "expecting: no error" )
  }

}
