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
  "bufio"
  "fmt"
  "github.com/SatoshiPortal/cam/errors"
  "github.com/SatoshiPortal/cam/globals"
  "github.com/SatoshiPortal/cam/output"
  "github.com/SatoshiPortal/cam/utils"
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
  return utils.SliceIndex( len(sourceList.Sources), func(i int) bool {
    return sourceList.Sources[i].String() == sourceString
  } )
}

func (sourceList *SourceList) GetSourceByIndex( index int ) ISource {
  if index < 0 || index >= len( sourceList.Sources ) {
    return nil
  }
  return sourceList.Sources[index]
}

func (sourceList *SourceList) GetSourceByString( sourceString string ) ISource {
  index := sourceList.SourceIndex( sourceString )
  if index == -1 {
    return nil
  }
  return sourceList.Sources[index]
}

func (sourceList *SourceList) GetSourceByHash( hash string ) ISource {
  index := utils.SliceIndex( len(sourceList.Sources), func(i int) bool {
    return sourceList.Sources[i].GetHash() == hash
  } )
  if index == -1 {
    return nil
  }
  return sourceList.Sources[index]
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

func (sourceList *SourceList) UpdateSource( sourceString string ) error {
  source := sourceList.GetSourceByString( sourceString )
  if source == nil {
    return errors.NO_SUCH_SOURCE
  }
  err := source.Update()
  if err != nil {
    return err
  }
  return nil
}

func (sourceList *SourceList) RemoveSource( sourceString string ) error {
  sourceIndex := sourceList.SourceIndex( sourceString )
  if sourceIndex == -1 {
    return errors.NO_SUCH_SOURCE
  }
  sourceList.Sources[sourceIndex].Cleanup()
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
