package main

type ShelfOptions struct {
	List             bool `short:"l" long:"list" desciption:"Print books in shelf order"`
	PrintAll         bool `short:"a" long:"all" description:"Print all books including those with LCC or Title"`
	NoInfo           bool `long:"no-info" description:"Only print books without information (ISBN only)"`
	NoClassification bool `long:"no-class" description:"Only print books without a LCC"`

	Resort bool `long:"resort" description:"Resort your library"`

	FindMissing bool `short:"f" long:"find-missing" description:"Search OpenLibrary for books missing a Title/LCC"`
	FindAll     bool `long:"find-all" description:"Search OpenLibrary for all books and update"`
}

func RunShelf(opts *ShelfOptions) {
	lib, err := LibraryLoad()
	if err != nil {
		Error(err.Error())
		return
	}

	opts.List = opts.List || opts.NoInfo || opts.NoClassification || opts.PrintAll

	if opts.FindMissing || opts.FindAll {
		shelfFind(lib, opts)
	}

	if opts.Resort {
		shelfResort(lib, opts)
	}
	if opts.List {
		shelfPrint(lib, opts)
	}
}

func shelfFind(lib *LibraryV1, opts *ShelfOptions) {
	for i, e := range lib.Entries {
		if opts.FindAll || e.LCC == "" || e.Title == "" {

			var ol *OpenLibraryBook
			var err error
			if e.OLID != "" {
				ol, err = OLGetByOLID(e.OLID)
			} else {
				ol, err = OLGetByISBN(e.FirstISBN())
			}

			if err != nil {
				Warn(err.Error())
				continue
			}
			lib.Entries[i].Update(ol)
		}
	}

	if err := LibrarySave(lib); err != nil {
		Error(err.Error())
	}
}

func shelfResort(lib *LibraryV1, opts *ShelfOptions) {
	err := LibrarySave(lib) // sorts on save
	if err != nil {
		Error(err.Error())
	}
}

func shelfPrint(lib *LibraryV1, opts *ShelfOptions) {
	var print func(LibraryEntryV1) bool
	if opts.PrintAll {
		print = func(e LibraryEntryV1) bool {
			if e.Title != "" {
				Print("%-20s - %s", e.LCC, e.Title)
			} else {
				Print("%-20s - ISBN: %s", e.LCC, e.FirstISBN())
			}
			return true
		}
	} else if opts.NoClassification {
		print = func(e LibraryEntryV1) bool {
			if e.LCC != "" {
				return false
			}
			if e.Title != "" {
				Print("%-50s https://openlibrary.org/books/%s", e.Title, e.OLID)
			} else {
				Print("ISBN: %-44s https://openlibrary.org/books/add", e.FirstISBN())
			}
			return true
		}
	} else if opts.NoInfo {
		print = func(e LibraryEntryV1) bool {
			if e.LCC != "" {
				return false
			}
			if e.Title != "" { // technically should check other fields, but Title is kind of the main thing to have
				return true
			}
			Print("%s", e.FirstISBN())
			return true
		}
	} else {
		print = func(e LibraryEntryV1) bool {
			if e.LCC != "" && e.Title != "" {
				Print("%-20s - %s", e.LCC, e.Title)
			}
			return true
		}
	}

	for _, e := range lib.Entries {
		if !print(e) {
			break
		}
	}
}
