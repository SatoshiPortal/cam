package cyphernode_test

import (
  "encoding/json"
  "github.com/schulterklopfer/cam/cyphernode"
  "testing"
)

const jsonString = `{
  "api_versions": [
    "v0"
  ],
  "bitcoin_version": "0.18.0",
  "optional_features": [
    {
      "label": "lightning",
      "active": true,
      "docker": "cyphernode/clightning:0.7.1"
    }
  ]
}`

func TestUnmarshalJSON(t *testing.T) {
  var cni cyphernode.CyphernodeInfo

  err := json.Unmarshal( []byte(jsonString), &cni )

  if err != nil {
    t.Error( err.Error() )
  }

}
