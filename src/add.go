package main

import (
	"fmt"
	"os"
)

type AddOptions struct {
	FilePath string `short:"f" long:"file" description:"File to load multiple books from"`
}

func RunAdd(opts *AddOptions) {
	lib, err := LibraryLoad()
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
	}
	_ = lib
}
