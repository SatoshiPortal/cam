package storage

import (
  "encoding/json"
  "github.com/schulterklopfer/cna/globals"
  "github.com/schulterklopfer/cna/utils"
  "github.com/schulterklopfer/cna/version"
  "path/filepath"
)

type App struct {
  Label string `json:"label"`
  Name string `json:"name"`
  Path string `json:"path"`
  URL string `json:"url"`
  Email string `json:"email"`
  Latest string `json:"latest"`
  Source ISource `json:"source"`
  Candidates []*AppCandidate `json:"candidates"`
  hash string `json:"-"`
}

type AppCandidate struct {
  Version *version.Version `json:"version"`
  Dependencies []*AppDependency `json:"dependencies"`
  Files []string `json:"files"`
}

type AppDependency struct {
  Label string `json:"label"`
  Version *version.Version `json:"version"`
}

/** App **/

func NewApp() *App {
  return &App{}
}

func (app *App) BuildHash() {
  bytes := make( []byte, 0 )
  bytes = append( bytes, []byte(app.Label)... )
  bytes = append( bytes, []byte(app.Source.GetHash())... )
  app.hash = utils.BuildHash( &bytes )
}

func (app *App) GetHash() string {
  return app.hash
}

/** AppCandidate **/

func NewAppCandidate () *AppCandidate {
  return &AppCandidate{}
}

func (appCandidate *AppCandidate) GetFilesDir() string {
  return filepath.Join(globals.APP_VERSIONS_DIR, appCandidate.Version.Raw)
}

/** AppDependency **/

func NewAppDependency () *AppDependency {
  return &AppDependency{}
}

func (app *App) UnmarshalJSON(data []byte) error {
  /** Big hack
    we need to to this, cause json.Unmarhsal is unable to assign
    the Source property because of the ISource interface
   **/

  intermediate := make( map[string]interface{} )
  err := json.Unmarshal( data, &intermediate )
  if err != nil {
    return err
  }

  for key, value := range intermediate {
    switch key {
    case "source":
      app.Source = SourceFromString( value.(string) )
      break
    case "candidates":
      candidatesJsonBytes, err := json.Marshal( value )
      if err != nil {
        return err
      }
      var candidates []*AppCandidate
      err = json.Unmarshal( candidatesJsonBytes, &candidates)
      if err != nil {
        return err
      }
      app.Candidates = candidates
      break
    case "label":
      app.Label = value.(string)
    case "name":
      app.Name = value.(string)
    case "path":
      app.Path = value.(string)
    case "url":
      app.URL = value.(string)
    case "email":
      app.Email = value.(string)
    case "latest":
      app.Latest = value.(string)
    }
  }
  return nil
}
