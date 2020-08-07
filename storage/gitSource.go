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
  "github.com/SatoshiPortal/cam/output"
  "github.com/SatoshiPortal/cam/utils"
  "gopkg.in/src-d/go-git.v4"
  "gopkg.in/src-d/go-git.v4/plumbing"
  "os"
  "strings"
)

type GitSource struct {
  Source
}

func NewGitSource( url string ) *GitSource {
  source := &GitSource{
    Source{location: url},
  }
  source.BuildHash()
  return source
}

func ( gitSource *GitSource ) BuildHash() {
  bytes := make( []byte, 0 )
  bytes = append( bytes, []byte(gitSource.location)... )
  gitSource.hash = utils.BuildHash( &bytes )
}

func ( gitSource *GitSource ) GetHash() string {
  return gitSource.hash
}

func ( gitSource *GitSource ) GetType() string {
  return SOURCE_TYPE_GIT
}

func ( gitSource *GitSource ) Update() error {
  url := strings.Replace( gitSource.location, "git://", "https://", 1)
  relativeRepoDir := strings.Replace( gitSource.location, "git://", "", 1)

  targetDir := utils.GetRepoDirPathFor( relativeRepoDir )

  var gitRepo *git.Repository
  var err error

  if repoDirExists, _ := utils.RepoExists( relativeRepoDir ); !repoDirExists {
    err = os.MkdirAll( targetDir, 0755)
    if err != nil {
      return err
    }
    // clone
    gitRepo, err = git.PlainClone(targetDir, false, &git.CloneOptions{
      URL: url,
      RemoteName: "origin",
      ReferenceName: plumbing.ReferenceName("refs/heads/cam"),
      Progress: os.Stdout,
    })
    if err != nil {
      return err
    }

  } else {
    // pull
    gitRepo, err = git.PlainOpen(targetDir)
    if err != nil {
      return err
    }
    // Get the working directory for the repository
    w, err := gitRepo.Worktree()
    if err != nil {
      return err
    }
    // Pull the latest changes from the origin remote and merge into the current branch
    err = w.Pull(&git.PullOptions{RemoteName: "origin"})
    if err != nil {
      output.Warning( err.Error() )
    }
  }

  // Print the latest commit that was just pulled
  ref, err := gitRepo.Head()
  if err != nil {
    return err
  }
  commit, err := gitRepo.CommitObject(ref.Hash())
  if err != nil {
    return err
  }

  output.Notice( commit.String() )

  return nil
}

func ( gitSource *GitSource ) String() string {
  return gitSource.location
}

func ( gitSource *GitSource ) GetAbsolutePath() string {
  return utils.GetRepoDirPathFor( strings.Replace( gitSource.location, "git://", "", 1) )
}

func (gitSource *GitSource) UnmarshalJSON(data []byte) error {
  var v string
  err := json.Unmarshal( data, &v )
  if err != nil {
    return err
  }
  gitSource.location = v
  return nil
}

func (gitSource *GitSource) MarshalJSON()  ([]byte, error)  {
  return json.Marshal( gitSource.location )
}

func (gitSource *GitSource) Cleanup() {
  dir := gitSource.GetAbsolutePath()
  err := os.RemoveAll(dir)
  if err != nil {
    output.Warningf("git source cleanup error: %s", err.Error() )
  }
}