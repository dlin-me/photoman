package main

import (
	"path/filepath"
	"io/ioutil"
	"encoding/gob"
	"bytes"
	"os"
)

func RestoreIndex(indexPath string) (map[string][]string, error) {
	decodedMap := make(map[string][]string)

	data, err := ioutil.ReadFile(indexPath)

	if err != nil {
		return decodedMap, err
	}

	b := bytes.NewBuffer(data)
	d := gob.NewDecoder(b)
	err = d.Decode(&decodedMap)

	return decodedMap, err
}

func RestoreHashIndex(dirPath string) (map[string][]string, error) {
	hashIndexPath := filepath.Join(dirPath, ".photoman_hdx")
	return RestoreIndex(hashIndexPath)
}


func RestorePathIndex(dirPath string) (map[string][]string, error) {
	hashIndexPath := filepath.Join(dirPath, ".photoman_pdx")
	return RestoreIndex(hashIndexPath)
}

func SaveIndex(index map[string][]string, indexPath string) error {
	file, err := os.Create(indexPath)

	if err != nil {
		return err
	}

	b := new(bytes.Buffer)
	e := gob.NewEncoder(b)
	err = e.Encode(index)

	if err != nil {
		return err
	}

	_, err = b.WriteTo(file)

	return err
}

func SaveHashIndex(hashMap map[string][]string, dirPath string) error {
	indexPath := filepath.Join(dirPath, ".photoman_hdx")
	return SaveIndex(hashMap, indexPath)
}

func SavePathIndex(pathMap map[string][]string, dirPath string) error {
	indexPath := filepath.Join(dirPath, ".photoman_pdx")
	return SaveIndex(pathMap, indexPath)
}

func IndexFile(filePath string, version string, hashMap map[string][]string, pathMap map[string][]string) error {
	_, ok := pathMap[filePath];

	if !ok {
		hash, err := ComputeMD5(filePath)

		if err != nil {
			return err
		}

		pathMap[filePath] = []string{hash, version}
		paths, ok := hashMap[hash]

		if ok {
			hashMap[hash] = append(paths, filePath)
		} else {
			hashMap[hash] = []string{filePath}
		}
	} else {
		pathMap[filePath][1] = version
	}

	return nil
}

func CleanIndex(versionToKeep string, hashMap map[string][]string, pathMap map[string][]string) error {
	var version string

	for path, data := range pathMap {
		_, version = data[0], data[1]

		if version != versionToKeep {
			RemoveFileFromIndex(path, hashMap, pathMap)
		}
	}

	return nil
}


func RemoveFileFromIndex(filePath string, hashMap map[string][]string, pathMap map[string][]string) {

	hashData, ok := pathMap[filePath];

	if ok {
		hash := hashData[0]

		paths, ok := hashMap[hash]

		if ok && len(paths) == 1 {
			delete(hashMap, hash)
		} else if ok && len(paths) > 1 {
			for i, path := range paths {
				if path == filePath {
					hashMap[hash] = append(hashMap[hash][:i], hashMap[hash][i+1:]...)
				}
			}
		}

		delete(pathMap, filePath)
	}
}

func ExtractDuplicatedFiles(hashMap map[string][]string) map[string][]string {
	result := make(map[string][]string)

	for hash, paths := range hashMap {
		if len(paths) > 1 {
			result[hash] = paths
		}
	}

	return result;
}





