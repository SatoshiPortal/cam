package storage

import (
  "github.com/schulterklopfer/cna/errors"
)

type AppList struct {
  Apps []*App `json:"data"`
  Labels map[string][]int `json:"labels"`
}

func (appList *AppList) AppIndex( app *App ) int {
  for i:=0; i<len( appList.Apps ); i++ {
    // TODO: fix
    if appList.Apps[i].GetHash() == app.GetHash() {
      return i
    }
  }
  return -1
}

func (appList *AppList) BuildLabels() {
  for i:=0; i<len( appList.Apps ); i++ {
    appList.buildLabel( appList.Apps[i] )
  }
}

func (appList *AppList) buildLabel( app *App ) {
  // 1) use app name for label

  // 2) use hash as label
}

func (appList *AppList) AddApp( app *App ) error {

  if appList.AppIndex( app ) >= 0 {
    return errors.DUPLICATE_APP
  }
  appList.Apps = append(appList.Apps, app)
  appList.buildLabel(app)
  return nil
}

func (appList *AppList) RemoveApp( app *App ) error {
  appIndex := appList.AppIndex( app )
  if appIndex == -1 {
    return errors.NO_SUCH_APP
  }
  appList.Apps = append(appList.Apps[:appIndex], appList.Apps[appIndex+1:]...)
  return nil
}

func (appList *AppList) Save() error {
  return nil
}
