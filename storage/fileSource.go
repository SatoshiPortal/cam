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

package storage

import (
  "encoding/json"
  "github.com/SatoshiPortal/cam/utils"
  "strings"
)

type FileSource struct {
  Source
}

func NewFileSource( url string ) *FileSource {
  source := &FileSource{
    Source{location: url},
  }
  source.BuildHash()
  return source
}

func ( fileSource *FileSource ) BuildHash() {
  bytes := make( []byte, 0 )
  bytes = append( bytes, []byte(fileSource.location)... )
  fileSource.hash = utils.BuildHash( &bytes )
}

func ( fileSource *FileSource ) GetHash() string {
  return fileSource.hash
}

func ( fileSource *FileSource ) GetType() string {
  return SOURCE_TYPE_FILE
}

func ( fileSource *FileSource ) Update() error {
  return nil
}

func ( fileSource *FileSource ) String() string {
  return fileSource.location
}

func ( fileSource *FileSource ) GetAbsolutePath() string {
  return strings.Replace( fileSource.location, "file://", "", 1)
}

func (fileSource *FileSource) UnmarshalJSON(data []byte) error {
  var v string
  err := json.Unmarshal( data, &v )
  if err != nil {
    return err
  }
  fileSource.location = v
  return nil
}

func (fileSource *FileSource) MarshalJSON()  ([]byte, error)  {
  return json.Marshal( fileSource.location )
}

func  (fileSource *FileSource) Cleanup() {

}