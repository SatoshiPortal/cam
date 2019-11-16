package storage

import (
  "github.com/schulterklopfer/cna/errors"
  "github.com/schulterklopfer/cna/utils"
  "strings"
)

type AppList struct {
  Apps []*App `json:"data,omitempty"`
  Labels map[string][]int `json:"-"`
}

func (appList *AppList) AppIndex( app *App ) int {
  return utils.SliceIndex( len(appList.Apps), func(i int) bool {
    return appList.Apps[i].GetHash() == app.GetHash()
  } )
}

func (appList *AppList) BuildLabels() {
  appList.Labels = make( map[string][]int )
  for i:=0; i<len( appList.Apps ); i++ {
    appList.buildLabel( i )
  }
}

func (appList *AppList) buildLabel( appIndex int ) {
  app := appList.Apps[appIndex]
  utils.AddIndexToLabel( &appList.Labels, app.Label, appIndex )
  utils.AddIndexToLabel( &appList.Labels, app.GetHash(), appIndex )
}

func (appList *AppList) AddApp( app *App ) error {

  if appList.AppIndex( app ) >= 0 {
    return errors.DUPLICATE_APP
  }
  appList.Apps = append(appList.Apps, app)
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

func (appList *AppList) Search( text string, exact bool ) []*App {
  apps := make( []*App, 0 )

  for label, indeces := range appList.Labels {

    found := false
    if exact {
      found = label==text
    } else {
      found = strings.Contains( label, text )
    }

    if found {
      for i:=0; i<len(indeces); i++ {
        apps = append( apps, appList.Apps[indeces[i]])
      }
    }
  }
  return apps
}

func (appList *AppList) BuildAppHashes() {
  for i:=0; i<len( appList.Apps ); i++ {
    appList.Apps[i].BuildHash()
  }
}

func (appList *AppList) Clear() {
  appList.Apps = appList.Apps[0:0]
  for k := range appList.Labels {
    delete(appList.Labels, k)
  }
}

func (appList *AppList) Save() error {
  return nil
}
