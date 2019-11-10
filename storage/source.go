package storage

const SOURCE_TYPE_GIT = "git"
const SOURCE_TYPE_FILE = "file"

type ISource interface {
  GetType() string
  Update() error
  String() string
  GetAbsolutePath() string
  GetHash() string
  BuildHash()
}

type Source struct {
  location string
  hash string
}

