package main

import (
	"path/filepath"
	"io/ioutil"
	"encoding/gob"
	"bytes"
	"os"
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

	return decodedMap, err
}

func SaveIndex(index map[string][]string, dirPath string) error {
	indexPath := filepath.Join(dirPath, ".photoman_db")

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

func IndexFile(filePath string, version string, dbMap map[string][]string) error {
	_, ok := dbMap[filePath];

	if !ok {
		hash, err := ComputeMD5(filePath)

		if err != nil {
			return err
		}

		dbMap[filePath] = []string{hash, version}

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





