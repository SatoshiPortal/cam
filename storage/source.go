package storage

import (
  "net/url"
)

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
  location string `json:"location"`
  hash string `json:"-"`
}

func SourceFromString( sourceString string ) ISource {
  sourceUrl, err := url.Parse( sourceString )
  if err != nil {
    return nil
  }
  switch sourceUrl.Scheme {
  case SOURCE_TYPE_GIT:
    return  NewGitSource( sourceString )
    break
  case SOURCE_TYPE_FILE:
    return NewFileSource( sourceString )
    break
  }
  return nil
}