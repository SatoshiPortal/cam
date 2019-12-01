package actions

import (
  "github.com/SatoshiPortal/cam/output"
  "github.com/SatoshiPortal/cam/storage"
  "github.com/SatoshiPortal/cam/utils"
  "github.com/urfave/cli"
)

func Global_init(c *cli.Context) error {
  err := storage.InitDataDir()
  if err != nil {
    return err
  }
  err = storage.InitInstallDir()
  if err != nil {
    return err
  }
  err = storage.Unlock()
  if err != nil {
    return err
  }
  return nil
}

func Global_update(c *cli.Context) error {
  sourceList, err := storage.LoadSourceFile( utils.GetSourceFilePath() )

  if err != nil {
    return err
  }

  for i:=0; i<len(sourceList.Sources); i++ {
    output.Noticef( "Updating %s\n", sourceList.Sources[i].String() )
    err = sourceList.Sources[i].Update()
    if err != nil {
      output.Errorf( "Error updating source: %s", sourceList.Sources[i].String() )
    }
  }


  repoIndex, err := storage.NewRepoIndex()

  if err == nil {
    output.Notice( "Recreating repo index")
  } else {
    output.Notice( "Building repo index")
  }

  err = repoIndex.Build()

  if err != nil {
    return err
  }

  return nil
}
