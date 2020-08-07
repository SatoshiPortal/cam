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