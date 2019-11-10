package storage

import (
  "github.com/schulterklopfer/cna/globals"
  "github.com/schulterklopfer/cna/utils"
  "github.com/schulterklopfer/cna/version"
  "path/filepath"
  "sort"
)

type App struct {
  Label string `json:"label"`
  Name string `json:"name"`
  Path string `json:"path"`
  URL string `json:"url"`
  Email string `json:"email"`
  Latest string `json:"latest"`
  Source ISource `json:"-"`
  Candidates []*AppCandidate `json:"candidates"`
  hash string `json:"-"`
}

type AppCandidate struct {
  Version *version.Version `json:"version"`
  Dependencies []*AppDependency `json:"dependencies"`
  Files []string `json:"files"`
  hash string `json:"-"`
}

type AppDependency struct {
  Label string `json:"label"`
  Version *version.Version `json:"version"`
  hash string `json:"-"`
}

/** App **/

func NewApp() *App {
  return &App{}
}

func (app *App) BuildHash() {
  candidates := app.Candidates[:]
  sort.SliceStable(candidates, func(i, j int) bool {
    return candidates[i].Version.Raw < candidates[j].Version.Raw
  })

  bytes := make( []byte, 0 )
  bytes = append( bytes, []byte(app.Label)... )
  bytes = append( bytes, []byte(app.Name)... )
  bytes = append( bytes, []byte(app.Source.GetHash())... )

  for i:=0; i< len(candidates); i++ {
    candidates[i].BuildHash()
    bytes = append( bytes, []byte(candidates[i].GetHash())... )
  }

  app.hash = utils.BuildHash( &bytes )
}

func (app *App) GetHash() string {
  return app.hash
}

/** AppCandidate **/

func NewAppCandidate () *AppCandidate {
  return &AppCandidate{}
}

func (appCandidate *AppCandidate) BuildHash() {
  dependencies := appCandidate.Dependencies[:]
  sort.SliceStable(dependencies, func(i, j int) bool {
    return dependencies[i].Label < dependencies[j].Label
  })

  bytes := make( []byte, 0 )
  bytes = append( bytes, []byte(appCandidate.Version.Raw)... )

  for i:=0; i< len(dependencies); i++ {
    dependencies[i].BuildHash()
    bytes = append( bytes, []byte(dependencies[i].GetHash())... )
  }

  appCandidate.hash = utils.BuildHash( &bytes )

}

func (appCandidate *AppCandidate) GetHash() string {
  return appCandidate.hash
}

func (appCandidate *AppCandidate) GetFilesDir() string {
  return filepath.Join(globals.APP_VERSIONS_DIR, appCandidate.Version.Raw)
}

/** AppDependency **/

func NewAppDependency () *AppDependency {
  return &AppDependency{}
}

func (appDependency *AppDependency) BuildHash() {
  bytes := make( []byte, 0 )
  bytes = append( bytes, []byte(appDependency.Label)... )
  bytes = append( bytes, []byte(appDependency.Version.Raw)... )
  appDependency.hash = utils.BuildHash( &bytes )
}

func (appDependency *AppDependency) GetHash() string {
  // TODO: build hash from sorted app keys
  return appDependency.hash
}
