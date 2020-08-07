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
  "github.com/SatoshiPortal/cam/globals"
  "io/ioutil"
  "time"
)

type State struct {
  Version string `json:"version"`
  LastUpdate time.Time `json:"lastUpdate"`
  Apps *AppList  `json:"apps"`
  InstalledApps *AppList `json:"installedApps"`
}

type StateFile struct {
  Path string
  State *State
}

func NewStateFile( path string ) *StateFile {
  return &StateFile{
    Path: path,
    State: &State{
      Version:    globals.VERSION,
      LastUpdate: time.Now(),
    },
  }
}

func LoadStateFile( path string ) (*StateFile, error) {
  stateJsonBytes, err := ioutil.ReadFile( path )
  if err != nil {
    return nil, err
  }

  var state State
  err = json.Unmarshal( stateJsonBytes, &state )
  if err != nil {
    return nil, err
  }

  return &StateFile{
    Path: path,
    State: &state,
  }, nil

}

func (stateFile *StateFile) Save() error {
  stateFile.State.LastUpdate = time.Now()
  stateJsonBytes, err := json.MarshalIndent( stateFile.State, "", "  " )
  err = ioutil.WriteFile(stateFile.Path, stateJsonBytes, 0644)
  if err != nil {
    return err
  }
  return nil
}
