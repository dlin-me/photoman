package main

import (
	// "fmt"
	"log"
	"os"
	"github.com/urfave/cli"
	// "github.com/cheggaaa/pb"
	// "path/filepath"

)

func panicIfErr(err error) {
	if err != nil {
		log.Panicln(err.Error())
	}
}

func main() {
	app := cli.NewApp()
  	app.Name = "Photo Manager"
  	app.Usage = "Organise your photo files"
	app.Version = "0.0.1"
	app.Commands = []cli.Command{
		IndexCommand,
		DeduplicateCommand,
		MoveCommand,
	}

	err := app.Run(os.Args)
	panicIfErr(err)
}
