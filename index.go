package main

import (
	"path/filepath"
	"io/ioutil"
	"encoding/gob"
	"bytes"
	"os"
	"errors"
)

func RestoreIndex(dirPath string) (map[string][]string, error) {
	indexPath := filepath.Join(dirPath, ".photoman_db")

	decodedMap := make(map[string][]string)

	data, err := ioutil.ReadFile(indexPath)

	if err != nil {
		return decodedMap, err
	}

	b := bytes.NewBuffer(data)
	d := gob.NewDecoder(b)
	err = d.Decode(&decodedMap)

	if err != nil {
		return decodedMap, err
	} else {
		absIndex := make(map[string][]string)

		for path, data := range decodedMap {
			absPath, err := filepath.Abs(path)

			if err != nil {
				return absIndex, err
			}

			absIndex[absPath] = data
		}

		return absIndex, err;
	}
}

func SaveIndex(index map[string][]string, dirPath string) error {
	indexPath := filepath.Join(dirPath, ".photoman_db")

	file, err := os.Create(indexPath)

	if err != nil {
		return err
	}

	// replace abs path with rel paths
	relIndex := make(map[string][]string)

	for path, data := range index {
		relPath, err := GetRelativePath(path)

		if err != nil {
			return err
		}

		relIndex[relPath] = data
	}

	b := new(bytes.Buffer)
	e := gob.NewEncoder(b)
	err = e.Encode(relIndex)

	if err != nil {
		return err
	}

	_, err = b.WriteTo(file)

	return err
}

func RenameFileInIndex(oldPath string, newPath string, dbMap map[string][]string) error {
	v, ok := dbMap[oldPath];

	if ok {
		dbMap[newPath] = v
		delete(dbMap, oldPath)
		return nil
	}else {
		return errors.New("Old file does not exist")
	}
}

func IndexFile(filePath string, version string, dbMap map[string][]string) error {
	_, ok := dbMap[filePath];

	if !ok {
		hash, err := ComputeMD5(filePath)

		if err != nil {
			return err
		}

		t, err := GetExifDateTime(filePath)

		if err == nil {
			dbMap[filePath] = []string{hash, version, t.Format("20060102150405")}
		} else {
			dbMap[filePath] = []string{hash, version}
		}

	} else {
		dbMap[filePath][1] = version
	}

	return nil
}

func CleanIndex(versionToKeep string, dbMap map[string][]string) error {
	var version string

	for path, data := range dbMap {
		_, version = data[0], data[1]

		if version != versionToKeep {
			delete(dbMap, path)
		}
	}

	return nil
}


func ExtractDuplicatedFiles(dbMap map[string][]string) map[string][]string {
	result := make(map[string][]string)

	for path, data := range dbMap {
		hash, _ := data[0], data[1]
		_, ok := result[hash]

		if ok {
			result[hash] = append(result[hash], path);
		}else {
			result[hash] = []string{path}
		}
	}

	for hash, paths := range result {
		if len(paths) <= 1 {
			delete(result, hash)
		}
	}

	return result;
}




