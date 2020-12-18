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

package actions

import (
  "github.com/SatoshiPortal/cam/output"
  "github.com/SatoshiPortal/cam/storage"
  "github.com/SatoshiPortal/cam/utils"
  "github.com/urfave/cli"
)

func Source_update(c *cli.Context) error {
  sourceList, err := storage.LoadSourceFile(  utils.GetSourceFilePath() )

  if err != nil {
    return err
  }

  var sources []storage.ISource

  if len(c.Args()) == 0 {
    // update all sources
    sources = sourceList.Sources
  } else {
    for _, arg := range c.Args() {
      // arg might be a hash
      source := sourceList.GetSourceByHash( arg )
      if source == nil {
        // not a hash
        // arg might be source string?
        source = sourceList.GetSourceByString( arg )
        if source == nil {
          output.Noticef("%s is not a valid source... skipping\n", arg)
          continue
        }
      }
      sources = append( sources, source )
    }
  }

  for _,source := range sources {
    err := source.Update()
    if err != nil {
      output.Noticef("Error updating source %s: %s\n", source.String(), err.Error() )
    } else {
      output.Noticef("Updated source %s\n", source.String() )
    }
  }

  repoIndex, err := storage.NewRepoIndex()

  if err != nil {
    return err
  }

  err = repoIndex.Build()

  if err != nil {
    return err
  }

  return nil
}
