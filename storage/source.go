package storage

const SOURCE_TYPE_GIT = "git"
const SOURCE_TYPE_FILE = "file"

type ISource interface {
  GetType() string
  Update() error
  String() string
}

type Source struct {
  location string
}

