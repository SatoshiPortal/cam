package main

import (
  "github.com/schulterklopfer/cam/commands"
  "github.com/schulterklopfer/cam/globals"
  "github.com/schulterklopfer/cam/output"
  "github.com/urfave/cli"
  "os"
)


func main() {

  var app = cli.NewApp()
  app.Name = globals.NAME
  app.Usage = ""
  app.Description = globals.DESCRIPTION
  app.Author = globals.AUTHOR
  app.Version = globals.VERSION
  app.Commands = []cli.Command{}

  commands.InitGlobalCommands( app )
  commands.InitAppCommands( app )
  commands.InitSourceCommands( app )
  commands.InitKeyCommands( app )

  err := app.Run(os.Args)
  if err != nil {
    output.Error( err.Error() )
  }
}
