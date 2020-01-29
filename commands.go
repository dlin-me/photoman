package main

import (
	"fmt"
	"strings"
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
	Action: func(c *cli.Context) error {
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
		counter := 0

		walkfunc := func(path string, fi os.FileInfo, err error) error {
			bar.Increment()
			counter++
			if counter % 2000 == 0 && !c.Bool("dry") {
				e = SaveIndex(dbIndex, dirPath)
				panicIfErr(e)
			}
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

		return nil
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
	Action: func(c *cli.Context) error {
		dirPath, e := os.Getwd()
		panicIfErr(e)

		dbIndex, _ := RestoreIndex(dirPath)
		duplicates := ExtractDuplicatedFiles(dbIndex)

		fmt.Printf("Found %v duplicated files\n", len(duplicates))

		if len(duplicates) == 0 {
			return nil
		}

		bar := pb.StartNew(len(duplicates))

		for _, paths := range duplicates {

			toKeep, e := FileToKeep(paths)
			panicIfErr(e)

			for _, path := range paths {

				if toKeep != path {
					if !c.Bool("dry") {
						e := os.Remove(path)
						panicIfErr(e)
						delete(dbIndex, path)
					} else {
						fmt.Printf("Duplicated file: %v\n", FormatDuplicatedFile(toKeep, path));
					}
				}
			}

			bar.Increment()
		}

		if !c.Bool("dry") {
			e = SaveIndex(dbIndex, dirPath)
			panicIfErr(e)
		}

		bar.FinishPrint("Deduplicated files removed.")

		return nil
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
		cli.BoolFlag{
			Name: "greedy,g",
			Usage: "Move files with modified datetime if exif data is not available",
		},
	},

	Action: func(c *cli.Context) error {
		if c.Bool("greedy") && len(c.Args()) < 1 {
			fmt.Println("You must specify a directory if you want to rename/move all files based on modified date")
			return nil
		}

		dirPath, e := os.Getwd()
		panicIfErr(e)

		dbIndex, _ := RestoreIndex(dirPath)
		toMove := make(map[string]string)

		for path, data := range dbIndex {
			var tm time.Time
			var e error

			relPath, e := GetRelativePath(path)
			panicIfErr(e)

			if len(data) == 3 {
				tm, e = time.Parse("20060102150405", data[2]);
				panicIfErr(e)

			} else if  c.Bool("greedy") && strings.HasPrefix(relPath, c.Args()[0] + string(os.PathSeparator)) {

				tm, e =  LastModTime(path)
				panicIfErr(e)
			}

			if !tm.IsZero() {
				proposedDir := GetProposedPath(tm)
				proposedPath := filepath.Join(proposedDir, filepath.Base(path))

				if proposedPath != path {
					toMove[path] = proposedPath
				}
			}
		}

		fmt.Printf("Found %v files for relocation\n", len(toMove))

		if len(toMove) == 0 {
			return nil
		}

		bar := pb.StartNew(len(toMove))
		counter := 0

		for path, proposedPath := range toMove {
			if !c.Bool("dry") {
				// normalise proposedPath
				postfix := 0
				for ExistFile(proposedPath) {
					postfix++
					proposedPath = PostfixFilePath(proposedPath, postfix)
				}

				// create dir
				e = os.MkdirAll(filepath.Dir(proposedPath), os.ModeDir|0750)
				panicIfErr(e)

				// move
				e = os.Rename(path, proposedPath)
				panicIfErr(e)

				// update index for moving
				RenameFileInIndex(path, proposedPath, dbIndex)

				// clear empty dir
				RemoveEmptyDir(filepath.Dir(path))

				counter++

				if counter % 1000 == 0 && !c.Bool("dry"){
					e = SaveIndex(dbIndex, dirPath)
					panicIfErr(e)
				}
			}

			bar.Increment()
		}

		if !c.Bool("dry") {
			e = SaveIndex(dbIndex, dirPath)
			panicIfErr(e)
		}

		bar.FinishPrint("Files relocated.")

		return nil
	},
}

