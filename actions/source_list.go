package actions

import (
  "github.com/SatoshiPortal/cam/output"
  "github.com/SatoshiPortal/cam/storage"
  "github.com/SatoshiPortal/cam/utils"
  "github.com/urfave/cli"
)

func Source_list(c *cli.Context) error {
  sourceList, err := storage.LoadSourceFile( utils.GetSourceFilePath() )

  if err != nil {
    return err
  }

  for i:=0; i<len(sourceList.Sources); i++ {
    output.Noticef( "%s (%24s)\n", sourceList.Sources[i].String(), sourceList.Sources[i].GetHash() )
  }

  return nil
}
