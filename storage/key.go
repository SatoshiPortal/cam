package storage

import (
  "strings"
)

type Key struct {
  Label string `json:"label"`
  Data string `json:"data"`
  Groups []string `json:"groups"`
}

//kapi_id="003";kapi_key="f0b8bb52f4c7007938757bcdfc73b452d6ce08cc0c660ce57c5464ae95f35417";kapi_groups="stats,watcher,spender,admin";eval ugroups_${kapi_id}=${kapi_groups};eval ukey_${kapi_id}=${kapi_key}
func NewKeyFromText( text string ) *Key {
  arr := strings.Split( text, ";" )
  key := &Key{ Groups: make( []string, 0 ) }

  for _, item := range arr {
    itemArr := strings.Split( strings.Trim(item, " \n"), "=" )
    if len(itemArr) != 2 {
      continue
    }
    k := strings.Trim( itemArr[0], " \"\n" )
    v := strings.Trim( itemArr[1], " \"\n" )

    switch strings.ToLower(k) {
    case "kapi_id":
      key.Label = v
      break
    case "kapi_key":
      key.Data = v
      break
    case "kapi_groups":
      groups := strings.Split(v, "," )
      for _, group := range groups {
        key.Groups = append( key.Groups, strings.Trim( group, " \n") )
      }
      break
    }
  }

  if key.Label == "" || key.Data == "" || len(key.Groups) == 0 {
    return nil
  }

  return key
}