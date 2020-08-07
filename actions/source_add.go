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
  "github.com/SatoshiPortal/cam/errors"
  "github.com/SatoshiPortal/cam/storage"
  "github.com/SatoshiPortal/cam/utils"
  "github.com/urfave/cli"
  "strings"
)

func Source_add(c *cli.Context) error {
  sourceList, err := storage.LoadSourceFile(  utils.GetSourceFilePath() )

  if err != nil {
    return err
  }

  if len(c.Args()) == 0 {
    return errors.SOURCE_ADD_NO_SOURCE
  }

  sourceString := strings.Trim( c.Args().Get(0), " \n")
  err = sourceList.AddSource( sourceString )
  if err != nil {
    return err
  }

  err = sourceList.Save()
  if err != nil {
    return err
  }

  return nil
}
