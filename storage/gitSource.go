package storage

import (
  "encoding/json"
  "github.com/schulterklopfer/cna/output"
  "github.com/schulterklopfer/cna/utils"
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
      ReferenceName: plumbing.ReferenceName("refs/heads/cna"),
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