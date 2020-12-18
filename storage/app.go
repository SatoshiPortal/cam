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
  "github.com/SatoshiPortal/cam/utils"
  "github.com/SatoshiPortal/cam/version"
  "path/filepath"
  "regexp"
  "strings"
)

/* DO NOT FORGET TO ADD PROPERTIES TO THE CUSTOM UNMARSHALLER (UnmarshalJSON)!!! */
type App struct {
  Label string `json:"label,omitempty"`
  Name string `json:"name,omitempty"`
  Path string `json:"path,omitempty"`
  URL string `json:"url,omitempty"`
  Email string `json:"email,omitempty"`
  Latest string `json:"latest,omitempty"`
  TrustZone string `json:"trustZone,omitempty"`
  DefaultMountPoint string `json:"defaultMountPoint,omitempty"`
  MountPoint string `json:"mountPoint,omitempty"`
  Source ISource `json:"source,omitempty"`
  Candidates []*AppCandidate `json:"candidates,omitempty"`
  hash string `json:"-"`
  ClientID string `json:"clientID,omitempty"`
  Secret string `json:"secret,omitempty"`
  Keys []*Key `json:"keys,omitempty"`
}

type AppCandidate struct {
  Version *version.Version `json:"version"`
  Dependencies []*AppDependency `json:"dependencies"`
  Files []string `json:"files"`
  AvailableRoles []*Role `json:"availableRoles"`
  AccessPolicies []*AccessPolicy `json:"accessPolicies"`
  IsExposed bool `json:"isExposed"`
  Port int `json:"port"`
}

type Role struct {
  Name string `json:"name"`
  Description string `json:"description"`
  AutoAssign bool `json:"autoAssign"`
}

type AccessPolicy struct {
  Roles []string `json:"roles"`
  Patterns []string `json:"resources"`
  Effect string `json:"effect"` // allow, deny
  Actions []string `json:"actions"` // get, post, delete, put, patch, options
}

func ( ap *AccessPolicy ) Check( method string, path string, roleNames []string ) bool {

  if method == "" {
    method = "*"
  }

  if ap.Effect == "" {
    ap.Effect = "deny"
  }

  trimmedLowercaseMethod := strings.ToLower(strings.Trim(method, " "))
  methodMatches := false
  for _, action := range ap.Actions {
    trimmedLowercaseAction := strings.ToLower(strings.Trim(action, " "))
    methodMatches = trimmedLowercaseAction == "*" || trimmedLowercaseMethod == trimmedLowercaseAction
    if methodMatches {
      break
    }
  }

  if !methodMatches  {
    if ap.Effect == "allow" {
      return false
    }
    return true
  }

  pathMatches := false
  for _, pattern := range ap.Patterns {
    var err error
    pathMatches, err = regexp.Match( pattern, []byte(path) )
    if err != nil {
      continue
    }

    if pathMatches {
      break
    }
  }

  if !pathMatches  {
    if ap.Effect == "allow" {
      return false
    }
  }

  roleMatches := false
  for _, requiredRoleName := range ap.Roles {
    if requiredRoleName == "*" {
      roleMatches = true
    } else {
      if roleNames == nil {
        continue
      }
      for _, roleName := range roleNames {
        if requiredRoleName==roleName {
          roleMatches = true
          break
        }
      }
    }

    if roleMatches {
      break
    }
  }

  if ap.Effect != "allow" {
    roleMatches = !roleMatches
  }

  return roleMatches
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
  if app.Label == "" || app.Source == nil {
    return
  }
  bytes := make( []byte, 0 )
  bytes = append( bytes, []byte(app.Label)... )
  bytes = append( bytes, []byte(app.Source.GetHash())... )
  app.hash = utils.BuildHash( &bytes )
}

func (app *App) GetHash() string {
  return app.hash
}

func (app *App) GetCandidate( version *version.Version ) *AppCandidate {
  for i:=0; i< len(app.Candidates); i++ {
    if app.Candidates[i].Version.IsEqual( version ) {
      return app.Candidates[i]
    }
  }
  return nil
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
    case "keys":
      keysJsonBytes, err := json.Marshal( value )
      if err != nil {
        return err
      }
      var keys []*Key
      err = json.Unmarshal( keysJsonBytes, &keys)
      if err != nil {
        return err
      }
      app.Keys = keys
      break
    case "label":
      app.Label = value.(string)
      break
    case "name":
      app.Name = value.(string)
      break
    case "path":
      app.Path = value.(string)
      break
    case "url":
      app.URL = value.(string)
      break
    case "email":
      app.Email = value.(string)
      break
    case "latest":
      app.Latest = value.(string)
      break
    case "clientID":
      app.ClientID = value.(string)
      break
    case "secret":
      app.Secret = value.(string)
      break
    case "mountPoint":
      app.MountPoint = value.(string)
      break
    case "defaultMountPoint":
      app.DefaultMountPoint = value.(string)
      break
    case "trustZone":
      app.TrustZone = value.(string)
      break
    }
  }
  return nil
}
