package main

import (
	"crypto/md5"
	"encoding/hex"
	"io"
	"os"
	"time"
	"path/filepath"
	"github.com/rwcarlsen/goexif/exif"
	"github.com/rwcarlsen/goexif/mknote"
	"strings"
	"errors"
	"strconv"
	"io/ioutil"
	"fmt"
	"regexp"
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

func FileToKeep(filePaths []string) (string, error) {
	oldestPath := filePaths[0]
	movedPathPattern := regexp.MustCompile(`[0-9]{4}_[0-9]{2}`)
	oldestModTime, err := LastModTime(filePaths[0])

	if err != nil {
		return oldestPath, err
	}

	for _, path := range filePaths[1:] {
		modTime, err := LastModTime(path)

		if err != nil {
			return oldestPath, err
		}

		if modTime.Before(oldestModTime) && !movedPathPattern.MatchString(oldestPath) {
			oldestPath = path
			oldestModTime = modTime
		}
	}

	return oldestPath, err
}

func GetExifDateTime(filePath string) (time.Time, error) {
	var tm time.Time

	if !strings.HasPrefix(strings.ToLower(filepath.Ext(filePath)), ".jpg") {
		return tm, errors.New("Type not supported")
	}

	f, err := os.Open(filePath)
	if err != nil {
		return tm, err
	}

	x, err := exif.Decode(f)
	if err != nil {
		return tm, err
	}

	return x.DateTime()
}

func FormatDuplicatedFile(path1 string, path2 string) string {
	byteArray1 := []byte(path1)
	byteArray2 := []byte(path2)
	splitPoint := 0

	for i, v := range byteArray1 {
		if v != byteArray2[i] {
			break
		}

		if v == '/' {
			splitPoint = i
		}
	}

	return string(byteArray1[0:splitPoint + 1]) + "[ " + string(byteArray1[splitPoint + 1:]) + " == " + string(byteArray2[splitPoint + 1:]) + " ]"
}

func GetProposedPath(t time.Time) string {
	path, _ := os.Getwd()
	return filepath.Join(
		path,
		string(os.PathSeparator),
		strconv.Itoa(t.Year()),
		string(os.PathSeparator),
		strconv.Itoa(t.Year()) + "_" + fmt.Sprintf("%.2d", int(t.Month())),
	)
}

func RemoveEmptyDir(path string) {
	if files, _ := ioutil.ReadDir(path); len(files) == 0 {
		os.Remove(path)

		RemoveEmptyDir(filepath.Dir(path))
	}
}

func ExistFile(path string) bool {
	if _, err := os.Stat(path); err == nil {
		return true;
	}

	return false;
}

func PostfixFilePath(path string, postfix int) string {
	extension := filepath.Ext(path)
	basePath := path[0:len(path)-len(extension)]

	return basePath + "_" + strconv.Itoa(postfix) + extension
}

func GetRelativePath(path string) (string, error) {
	currentPath, e := os.Getwd()

	if e != nil {
		return "", e
	}

	return filepath.Rel(currentPath, path)
}
