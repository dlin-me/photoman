package main

import (
	"fmt"
	"github.com/codegangsta/cli"
	"github.com/cheggaaa/pb"
	"os"
	"time"
	"path/filepath"
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
		dirPath, e := os.Getwd()
		panicIfErr(e)

		fmt.Println("Scanning files, this may take a few minutes...")
		count, e := CountFilesInPath(dirPath)
		panicIfErr(e)

		dbIndex, e := RestoreIndex(dirPath)

		if len(dbIndex) < 0 {
			fmt.Printf("Indexing %v files...\n", count)
		} else {
			fmt.Printf("Updating index for %v files...\n", count)
		}

		t := time.Now()
		version := t.Format("20060102150405")

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

					fmt.Printf("Duplicated file: %v\n", FormatDuplicatedFile(toKeep, path));
				}
			}
		}

		if !c.Bool("dry") {
			e = SaveIndex(dbIndex, dirPath)
			panicIfErr(e)
		}
	},
}

var MoveCommand = cli.Command{
	Name:        "move",
	Usage:       "Move files with exif data to proposed directory YYYY/YYYY_MM",
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
		counter := 0

		for path, data := range dbIndex {

			if len(data) == 3 {
				tm, e := time.Parse("20060102150405", data[2]);
				panicIfErr(e)
				proposedDir := GetProposedPath(tm)
				proposedPath := filepath.Join(proposedDir, filepath.Base(path))

				if proposedPath != path {
					if !c.Bool("dry") {
						// move files now
						os.MkdirAll(proposedDir, 0777)
						e = os.Rename(path, proposedPath)
						panicIfErr(e)

						// update index for moving
						RenameFileInIndex(path, proposedPath, dbIndex)
					}

					counter++
				}
			}
		}

		fmt.Printf("Moved %v files\n", counter)
	},
}

