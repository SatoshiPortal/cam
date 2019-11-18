package storage

import (
  "encoding/json"
  goErrors "github.com/pkg/errors"
  "github.com/schulterklopfer/cam/errors"
  "github.com/schulterklopfer/cam/utils"
  "github.com/schulterklopfer/cam/version"
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

func AppCandidateIsRunnableOnCyphernode( appCandidate *AppCandidate ) (bool, error) {

  if !utils.CyphernodeInfoFileExists() {
    return false, errors.CYPHERNODE_INFO_FILE_DOES_NOT_EXIST
  }

  var cyphernodeInfo CyphernodeInfo

  cyphernodeInfoJsonBytes, err := ioutil.ReadFile( utils.GetCyphernodeInfoFilePath() )
  if err != nil {
    return false, err
  }

  err = json.Unmarshal( cyphernodeInfoJsonBytes, &cyphernodeInfo )
  if err != nil {
    return false, err
  }

  for _, dependency := range appCandidate.Dependencies {

    if dependency.Label == "api" {
      if utils.SliceIndex( len(cyphernodeInfo.ApiVersions), func(i int) bool {
        return cyphernodeInfo.ApiVersions[i] == dependency.Version.Raw
      } ) == -1 {
        return false, nil
      }
    } else {
      cyphernodeFeature := findCyphernodeFeature( &cyphernodeInfo, dependency.Label )
      if cyphernodeFeature == nil {
        return false, nil
      }

      if !dependency.Version.IsCompatible(cyphernodeFeature.Docker.Version) {
        return false, nil
      }
    }
  }

  return true, nil
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