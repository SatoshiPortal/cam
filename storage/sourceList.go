package storage

import (
  "bufio"
  "fmt"
  "github.com/schulterklopfer/cna/errors"
  "github.com/schulterklopfer/cna/globals"
  "github.com/schulterklopfer/cna/output"
  "os"
  "strings"
)

type SourceList struct {
  Path string
  Sources []ISource
}

func NewSourceList( path string ) *SourceList {
  return &SourceList{
    Path: path,
    Sources: []ISource{
      NewGitSource( globals.CYPHERAPPS_REPO ),
    },
  }
}

func LoadSourceFile( path string ) (*SourceList, error) {
  sourceFile, err := os.Open(path)
  if err != nil {
    return nil, err
  }
  defer sourceFile.Close()

  r := &SourceList{
    Path: path,
  }

  scanner := bufio.NewScanner(sourceFile)
  for scanner.Scan() {
    text := scanner.Text()
    if strings.HasPrefix( text, "#" ) {
      continue
    }
    err = r.AddSource(text)
    if err != nil {
      output.Warning( err.Error() )
    }
  }

  if err := scanner.Err(); err != nil {
    return nil, err
  }

  return r, nil
}

func (sourceList *SourceList) SourceIndex( sourceString string ) int {
  for i:=0; i<len( sourceList.Sources ); i++ {
    if sourceList.Sources[i].String() == sourceString {
      return i
    }
  }
  return -1
}

func (sourceList *SourceList) AddSource( sourceString string ) error {

  if sourceList.SourceIndex( sourceString ) >= 0 {
    return errors.DUPLICATE_SOURCE
  }

  source := SourceFromString( sourceString )
  if source != nil {
    sourceList.Sources = append( sourceList.Sources, source )
  }
  return nil
}

func (sourceList *SourceList) RemoveSource( sourceString string ) error {
  sourceIndex := sourceList.SourceIndex( sourceString )
  if sourceIndex == -1 {
    return errors.NO_SUCH_SOURCE
  }
  sourceList.Sources = append(sourceList.Sources[:sourceIndex], sourceList.Sources[sourceIndex+1:]...)
  return nil
}

func (sourceList *SourceList) Save() error {

  f, err := os.Create(sourceList.Path)
  defer f.Close()

  if err != nil {
    return err
  }

  for i:=0; i<len(sourceList.Sources); i++ {
    source := sourceList.Sources[i]
    _, err := fmt.Fprintf( f, "%s\n", source.String() )
    if err != nil {
      continue
    }
  }

  return nil
}
