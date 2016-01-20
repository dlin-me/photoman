package main

import (	
	"os"
	"path/filepath"
)

// index files in a given path and store index in a LevelDB.
func WalkFilesInPath(dirPath string, callback filepath.WalkFunc) error {
    
    fullPath, err := filepath.Abs(dirPath)

    if err != nil {
        return err
    }

	cb := func(path string, fi os.FileInfo, err error) error {        

        goodFile, res := TestFile(fi)

        if !goodFile {
            return res
        }

		return callback(path, fi, err)
	}
    
    return filepath.Walk(fullPath, cb)    
}

func CountFilesInPath(dirPath string) (int, error) {
    counter := 0
    cb := func(path string, fi os.FileInfo, err error) error {      
        counter++
        return nil
    }

    return counter, WalkFilesInPath(dirPath, cb)
}
  