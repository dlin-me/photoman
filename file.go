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

func OldestFile(filePaths []string) (string, error) {
	oldestPath := filePaths[0]
	oldestModTime, err := LastModTime(filePaths[0])

	if err != nil {
		return oldestPath, err
	}

	for _, path := range filePaths[1:] {
		modTime, err := LastModTime(path)

		if err != nil {
			return oldestPath, err
		}

		if modTime.Before(oldestModTime) {
			oldestPath = path
			oldestModTime = modTime
		}
	}

	return oldestPath, err
}
