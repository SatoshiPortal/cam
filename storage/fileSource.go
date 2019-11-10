package storage

import (
  "github.com/schulterklopfer/cna/utils"
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
  bytes = append( bytes, []byte(SOURCE_TYPE_FILE)... )
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