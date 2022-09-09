package main

import (
	"fmt"
	"os"
)

type InitOptions struct {
}

func RunInit(opts *InitOptions) {
	if LibraryExists() {
		fmt.Fprintf(os.Stderr, "Library file %s already exists\n", LibraryPath())
	}

	if err := LibrarySave(&LibraryV1{SortedBy: SortLCC}); err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
	}
}
