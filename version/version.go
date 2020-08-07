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

package version

import (
  "encoding/json"
  "github.com/pkg/errors"
  "strconv"
  "strings"
)

type Version struct {
  Prefix string
  Major int
  Minor int
  Patch int
  Raw string
  Misc string
  Semantic bool
}

func NewVersion( raw string ) *Version {
  version := &Version{}
  _ = version.Parse(raw)
  return version
}

func (version *Version) Parse( raw string ) error {
  prefix := ""
  version.Raw = raw

  if strings.HasPrefix(raw, "~" ) || strings.HasPrefix(raw, "^" ) {
    // ^ fix major
    // ~ fix minor
    prefix = raw[:1]
    raw = raw[1:]
  }

  // handle versions like v0.1.2
  if strings.HasPrefix(raw, "v" ) {
    raw = raw[1:]
  }

  vArr := strings.Split( raw, "." )
  vArrInt := make( []int, 3 )

  if len(vArr) < 3 {
    return errors.New( "Unknown version format: "+raw )
  }

  misc := ""
  if len(vArr) > 3 {
    misc = strings.Join( vArr[3:], "." )
    vArr = vArr[0:3]
  }

  miscIndex := strings.IndexAny( vArr[2], "-_" )

  if miscIndex != -1 {
    m :=  vArr[2][miscIndex+1:]
    vArr[2] = vArr[2][0:miscIndex]
    if misc == "" {
      misc = m
    } else {
      misc = m+"."+misc
    }
  }

  for i:=0; i<3; i++ {
    v, err := strconv.Atoi(vArr[i])
    if err != nil {
      return errors.New("Unknown version format: "+raw )
    }
    vArrInt[i] = v
  }

  version.Prefix = prefix
  version.Major = vArrInt[0]
  version.Minor = vArrInt[1]
  version.Patch = vArrInt[2]
  version.Misc = misc
  version.Semantic = true

  return nil
}

func (version *Version) IsCompatible( otherVersion *Version ) bool {

  // ignore other versions prefix
  if version.Prefix == "^" {
    // fixed major
    return version.Major == otherVersion.Major
  } else if version.Prefix == "~" {
    // fixed minor
    return version.Major == otherVersion.Major &&
        version.Minor == otherVersion.Minor
  }
  // exact match
  return version.Major == otherVersion.Major &&
      version.Minor == otherVersion.Minor &&
      version.Patch == otherVersion.Patch
}

func (version *Version) IsEqual( otherVersion *Version ) bool {
  return version.Major == otherVersion.Major &&
      version.Minor == otherVersion.Minor &&
      version.Patch == otherVersion.Patch &&
      version.Misc == otherVersion.Misc
}

func (version *Version) UnmarshalJSON(data []byte) error {
  var v string
  err := json.Unmarshal( data, &v )
  if err != nil {
    return err
  }
  _ = version.Parse( v )
  return nil
}

func (version *Version) MarshalJSON()  ([]byte, error)  {
  return json.Marshal( version.Raw )
}