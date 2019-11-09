package storage

type FileSource struct {
  Source
}

func NewFileSource( url string ) *FileSource {
  return &FileSource{
    Source{location: url},
  }
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