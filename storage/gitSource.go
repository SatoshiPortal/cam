package storage

import (
  "github.com/schulterklopfer/cna/output"
  "github.com/schulterklopfer/cna/utils"
  "gopkg.in/src-d/go-git.v4"
  "os"
  "strings"
)

type GitSource struct {
  Source
}

func NewGitSource( url string ) *GitSource {
  return &GitSource{
    Source{location: url},
  }
}

func ( gitSource *GitSource ) GetType() string {
  return SOURCE_TYPE_GIT
}

func ( gitSource *GitSource ) Update() error {
  url := strings.Replace( gitSource.location, "git://", "https://", 1)
  relativeRepoDir := strings.Replace( gitSource.location, "git://", "", 1)

  targetDir, err := utils.GetRepoDirPathFor( relativeRepoDir )
  if err != nil {
    return err
  }

  var gitRepo *git.Repository

  if repoDirExists, _ := utils.RepoExists( relativeRepoDir ); !repoDirExists {
    err = os.MkdirAll( targetDir, 0755)
    if err != nil {
      return err
    }
    // clone
    gitRepo, err = git.PlainClone(targetDir, false, &git.CloneOptions{
      URL: url,
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

  if err != nil {
    return err
  }

  return nil
}

func ( gitSource *GitSource ) String() string {
  return gitSource.location
}