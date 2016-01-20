package main

import (
	"fmt"
	"github.com/codegangsta/cli"
	"github.com/cheggaaa/pb"
	"os"
	"time"
)

var IndexCommand = cli.Command{
	Name:        "index",
	Usage:       "Build an index for files in the current directory",
	Flags: []cli.Flag {
		cli.BoolFlag{
			Name: "dry,d",
			Usage: "Dry run only",
		},
	},
	Action: func(c *cli.Context) {
		fmt.Println("Counting files, this may take a few minutes...")

		dirPath, e := os.Getwd()
		panicIfErr(e)

		count, e := CountFilesInPath(dirPath)
		panicIfErr(e)

		fmt.Printf("Indexing %v files...\n", count)

		t := time.Now()
		version := t.Format("20060102150405")

		dbIndex, _ := RestoreIndex(dirPath)

 		bar := pb.StartNew(count)

		walkfunc := func(path string, fi os.FileInfo, err error) error {
			bar.Increment()
			return IndexFile(path, version, dbIndex)
		}

		e = WalkFilesInPath(dirPath, walkfunc)
		panicIfErr(e)

		CleanIndex(version, dbIndex)

		if !c.Bool("dry") {
			e = SaveIndex(dbIndex, dirPath)
			panicIfErr(e)
		}

		bar.FinishPrint("Index completed.")
	},
}

var DeduplicateCommand = cli.Command{
	Name:        "dd",
	Usage:       "Remove duplicated files",
	Flags: []cli.Flag {
		cli.BoolFlag{
			Name: "dry,d",
			Usage: "Dry run only",
		},
	},
	Action: func(c *cli.Context) {
		dirPath, e := os.Getwd()
		panicIfErr(e)

		dbIndex, _ := RestoreIndex(dirPath)
		duplicates := ExtractDuplicatedFiles(dbIndex)

		fmt.Printf("Found %v duplicated files\n", len(duplicates))

		for _, paths := range duplicates {

			toKeep, e := OldestFile(paths)

			panicIfErr(e)

			for _, path := range paths {

				if toKeep != path {

					if  !c.Bool("dry") {
						e := os.Remove(path)
						panicIfErr(e)
						delete(dbIndex, path)
					}

					fmt.Printf("Duplicated file: %v\n", path);
				}
			}
		}

		if !c.Bool("dry") {
			e = SaveIndex(dbIndex, dirPath)
			panicIfErr(e)
		}
	},
}

