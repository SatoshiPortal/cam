package storage

import (
  "bufio"
  "github.com/SatoshiPortal/cam/errors"
  "github.com/SatoshiPortal/cam/utils"
  "os"
  "strings"
)

type KeyList struct {
  Keys []*Key `json:"data,omitempty"`
  Labels map[string]int `json:"-"`
}

func NewKeyList() *KeyList {
  return &KeyList{Keys:make([]*Key, 0)}
}

func (keyList *KeyList) KeyIndex( key *Key ) int {
  return utils.SliceIndex( len(keyList.Keys), func(i int) bool {
    return keyList.Keys[i].Label == key.Label
  } )
}

func (keyList *KeyList) BuildLabels() {
  keyList.Labels = make( map[string]int )
  for i:=0; i<len( keyList.Keys ); i++ {
    keyList.Labels[keyList.Keys[i].Label] = i
  }
}

func (keyList *KeyList) AddKey( key *Key ) error {

  if keyList.KeyIndex(key) >= 0 {
    return errors.DUPLICATE_KEY
  }
  keyList.Keys = append(keyList.Keys, key)
  return nil
}

func (keyList *KeyList) RemoveKey( key *Key ) error {
  keyIndex := keyList.KeyIndex(key)
  if keyIndex == -1 {
    return errors.NO_SUCH_KEY
  }
  keyList.Keys = append(keyList.Keys[:keyIndex], keyList.Keys[keyIndex+1:]...)
  return nil
}

func (keyList *KeyList) Load() error {
  keysFilePath := utils.GetKeysFilePath()

  keysListFile, err := os.Open(keysFilePath)
  if err != nil {
    return err
  }
  defer keysListFile.Close()

  scanner := bufio.NewScanner(keysListFile)
  for scanner.Scan() {
    text := scanner.Text()
    if strings.HasPrefix( text, "#" ) {
      continue
    }
    // Parse keys line here
    key := NewKeyFromText( text )
    if key != nil {
      keyList.AddKey( key )
    }
  }

  if err := scanner.Err(); err != nil {
    return err
  }
  keyList.BuildLabels()
  return nil
}

func (keyList *KeyList) GetKey( label string ) *Key {
  if keyIndex, ok := keyList.Labels[label]; ok {
    return keyList.Keys[keyIndex]
  }
  return nil
}

func (keyList *KeyList) HasKey( label string ) bool {
  if _, ok := keyList.Labels[label]; ok {
    return ok
  }
  return false
}
