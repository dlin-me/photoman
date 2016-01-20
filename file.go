package main

import (
	"crypto/md5"
	"encoding/hex"
	"io"
	"os"
	"time"
	"path/filepath"
)

// computes MD5 hash for given file.
func ComputeMD5(filePath string) (string, error) {
	var result string
	file, err := os.Open(filePath)

	if err != nil {
		return result, err
	}

	defer file.Close()

	hash := md5.New()

	if _, err := io.Copy(hash, file); err != nil {
		return result, err
	}

	result = hex.EncodeToString(hash.Sum(nil))
	return result, nil
}


func LastModTime(filePath string) (time.Time, error) {
	info, err := os.Stat(filePath)

	if err != nil {
		return time.Now(), err
	}

	return info.ModTime(), err
}

func TestFile(fi os.FileInfo) (bool, error) {
	// Skip directories
	if fi.IsDir() {
		if filepath.HasPrefix(fi.Name(), ".") {
			return false, filepath.SkipDir
		} else {
			return false, nil
		}
	}

	// Skip hidden files
	if filepath.HasPrefix(fi.Name(), ".") {
		return false, nil
	}

	// Skip symlink
	if fi.Mode() & os.ModeSymlink == os.ModeSymlink {
		return false, nil
	}

	return true, nil
}

func KeepOldestFile(filePaths []string) ([]string, error) {

	oldestPath := filePaths[0]
	oldestModTime, err := LastModTime(filePaths[0])
	fileRemoved := []string{}

	if err != nil {
		return fileRemoved, err
	}

	for _, path := range filePaths[1:] {
		modTime, err := LastModTime(path)

		if err != nil {
			return fileRemoved, err
		}

		if modTime.Before(oldestModTime) {
			err = os.Remove(oldestPath)
			fileRemoved = append(fileRemoved, oldestPath)
			oldestPath = path
			oldestModTime = modTime
		} else {
			err = os.Remove(path)
			fileRemoved = append(fileRemoved, path)
		}

		if err != nil {
			return fileRemoved, err
		}

	}

	return fileRemoved, nil
}
