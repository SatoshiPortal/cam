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

package storage_test

import (
  "github.com/SatoshiPortal/cam/globals"
  "github.com/SatoshiPortal/cam/storage"
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
