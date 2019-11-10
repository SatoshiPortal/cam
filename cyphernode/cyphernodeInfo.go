package cyphernode

import (
  "encoding/json"
  "github.com/pkg/errors"
  "github.com/schulterklopfer/cna/version"
  "strings"
)

type CyphernodeInfo struct {
  ApiVersions []string `json:"api_versions"`
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
  v := strings.Trim( string(data), "\"")

  arr := strings.Split( v, ":" )
  if len(arr) == 0 {
    return errors.New( "unknown docker image format" )
  }

  dockerImage.ImageName = arr[0]
  if len(arr) > 1 {
    dockerImage.Version = version.NewVersion( arr[1] )
  }

  if err := json.Unmarshal(data, &v); err != nil {
    return err
  }
  return nil
}