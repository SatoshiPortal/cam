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
  goErrors "github.com/pkg/errors"
  "github.com/SatoshiPortal/cam/errors"
  "github.com/SatoshiPortal/cam/utils"
  "github.com/SatoshiPortal/cam/version"
  "io/ioutil"
  "strings"
)

type CyphernodeInfo struct {
  ApiVersions []string `json:"api_versions"`
  Features []*CyphernodeFeature `json:"features"`
  OptionalFeatures []*CyphernodeFeature `json:"optional_features"`
  BitcoinCoreVersion *version.Version `json:"bitcoin_version"`
}

type CyphernodeFeature struct {
  Label string `json:"label"`
  Active bool `json:"active"`
  Docker *DockerImage `json:"docker"`
}

type DockerImage struct {
  ImageName string
  Version *version.Version
}

func (dockerImage *DockerImage) UnmarshalJSON(data []byte) error {
  var v string

  if err := json.Unmarshal(data, &v); err != nil {
    return err
  }

  arr := strings.Split( v, ":" )
  if len(arr) == 0 {
    return goErrors.New( "unknown docker image format" )
  }

  dockerImage.ImageName = arr[0]
  if len(arr) > 1 {
    dockerImage.Version = version.NewVersion( arr[1] )
  }

  return nil
}

func AppCandidateIsRunnableOnCyphernode( appCandidate *AppCandidate ) error {

  if !utils.CyphernodeInfoFileExists() {
    return errors.CYPHERNODE_INFO_FILE_DOES_NOT_EXIST
  }

  var cyphernodeInfo CyphernodeInfo

  cyphernodeInfoJsonBytes, err := ioutil.ReadFile( utils.GetCyphernodeInfoFilePath() )
  if err != nil {
    return err
  }

  err = json.Unmarshal( cyphernodeInfoJsonBytes, &cyphernodeInfo )
  if err != nil {
    return err
  }

  for _, dependency := range appCandidate.Dependencies {

    if dependency.Label == "api" {
      if utils.SliceIndex( len(cyphernodeInfo.ApiVersions), func(i int) bool {
        return cyphernodeInfo.ApiVersions[i] == dependency.Version.Raw
      } ) == -1 {
        return errors.COMPAT_API
      }
    } else {
      cyphernodeFeature := findCyphernodeFeature( &cyphernodeInfo, dependency.Label )
      if cyphernodeFeature == nil || !cyphernodeFeature.Active {
        return errors.COMPAT_MISSING_FEATURE
      }

      if !dependency.Version.IsCompatible(cyphernodeFeature.Docker.Version) {
        return errors.COMPAT_FEATURE_VERSION
      }
    }
  }

  return nil
}

func findCyphernodeFeature( cyphernodeInfo *CyphernodeInfo, label string ) *CyphernodeFeature {
  for _, feature := range cyphernodeInfo.Features {
    if feature.Label == label {
      return feature
    }
  }
  for _, feature := range cyphernodeInfo.OptionalFeatures {
    if feature.Label == label {
      return feature
    }
  }
  return nil
}