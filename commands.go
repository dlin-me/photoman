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
	Action: func(c *cli.Context) {
		fmt.Println("Counting files, this may take a few minutes...")

		dirPath, e := os.Getwd()
		panicIfErr(e)

		count, e := CountFilesInPath(dirPath)
		panicIfErr(e)

		fmt.Printf("Indexing %v files...\n", count)

		t := time.Now()
		version := t.Format("20060102150405")

		hashIndex, _ := RestoreHashIndex(dirPath)
		pathIndex, _ := RestorePathIndex(dirPath)

 		bar := pb.StartNew(count)

		walkfunc := func(path string, fi os.FileInfo, err error) error {
			bar.Increment()
			return IndexFile(path, version, hashIndex, pathIndex)
		}

		e = WalkFilesInPath(dirPath, walkfunc)
		panicIfErr(e)

		CleanIndex(version, hashIndex, pathIndex)

		e = SaveHashIndex(hashIndex, dirPath)
		panicIfErr(e)
		e = SavePathIndex(pathIndex, dirPath)
		panicIfErr(e)

		bar.FinishPrint("Index completed.")
	},
}

var DeduplicateCommand = cli.Command{
	Name:        "dd",
	Usage:       "Remove duplicated files",
	Action: func(c *cli.Context) {
		dirPath, e := os.Getwd()
		panicIfErr(e)

		hashIndex, _ := RestoreHashIndex(dirPath)
		pathIndex, _ := RestorePathIndex(dirPath)

		duplicates := ExtractDuplicatedFiles(hashIndex)

		fmt.Printf("Found %v duplicated files\n", len(duplicates))

		for _, paths := range duplicates {

			fileRemoved, e := KeepOldestFile(paths)

			panicIfErr(e)

			for _, removedPath := range fileRemoved {
				RemoveFileFromIndex(removedPath, hashIndex, pathIndex)

				fmt.Printf("Removed file: %v\n", removedPath);
			}
		}

		e = SaveHashIndex(hashIndex, dirPath)
		panicIfErr(e)
		e = SavePathIndex(pathIndex, dirPath)
		panicIfErr(e)
	},
}

