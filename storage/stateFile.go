package storage

import (
  "encoding/json"
  "github.com/schulterklopfer/cna/globals"
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
