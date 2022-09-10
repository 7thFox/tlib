package main

type InitOptions struct {
}

func RunInit(opts *InitOptions) {
	if LibraryExists() {
		Error("Library file %s already exists", LibraryPath())
	}

	if err := LibrarySave(&LibraryV1{SortedBy: SortLCC}); err != nil {
		Error(err.Error())
	}
}
