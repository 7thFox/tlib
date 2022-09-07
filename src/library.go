package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
)

type Library struct {
	V1 *LibraryV1 `json:"v1"`
}

type LibraryV1 struct {
}

func LibraryLoad() (*LibraryV1, error) {
	Debug("Loading library from %s", LibraryPath())

	if !LibraryExists() {
		return nil, fmt.Errorf("Library file %s does not exist", LibraryPath())
	}

	bytes, err := ioutil.ReadFile(LibraryPath())
	if err != nil {
		return nil, err
	}

	var lib Library
	err = json.Unmarshal(bytes, &lib)
	if err != nil {
		return nil, err
	}

	if lib.V1 == nil {
		return nil, errors.New("V1 library not available")
	}

	return lib.V1, nil
}

func LibraryExists() bool {
	stat, err := os.Stat(LibraryPath())
	return err == nil && !stat.IsDir()
}

func LibraryPath() string {
	path := GlobalOpts.FilePath
	if path != "" && path[len(path)-1] != '/' {
		path += "/"
	}
	return path + "tlib.json"
}

func LibrarySave(lib *LibraryV1) error {
	data, err := json.Marshal(Library{
		V1: lib,
	})

	if err != nil {
		return err
	}

	return ioutil.WriteFile(LibraryPath(), data, 0644)
}
