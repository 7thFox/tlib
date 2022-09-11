package main

import (
	"bufio"
	"os"
	"strconv"
	"strings"
)

type AddOptions struct {
	FilePath        string `short:"f" long:"file" description:"File to load multiple books from"`
	NoImpliedISBN10 bool   `long:"no-implied-isbn10" description:"Disables implied leading 0 for ISBN's with only 9 characters"`
	StopOnError     bool   `short:"s" long:"stop-on-error" description:"Stops file add after first error"`
	UseStdin        bool   `long:"stdin" description:"Read from Stdin"`

	ByOLID bool `short:"o" long:"olid" description:"Import using OpenLibrary ID (eg OL35693885M)"`

	Positional struct {
		ISBN string
	} `positional-args:"true"`
}

func RunAdd(opts *AddOptions) {
	lib, err := LibraryLoad()
	if err != nil {
		Error(err.Error())
		return
	}

	isbns := make(map[string]bool, len(lib.Entries))
	for _, e := range lib.Entries {
		for _, i := range e.ISBN10 {
			isbns[i] = true
		}
		for _, i := range e.ISBN13 {
			isbns[i] = true
		}
	}

	var getBook func(string, *AddOptions) *LibraryEntryV1
	if opts.ByOLID {
		getBook = getBookByOLID
	} else {
		getBook = getBookByISBN
	}

	addBook := func(book *LibraryEntryV1, originalISBN string) bool {
		if book != nil {
			for _, i := range book.ISBN10 {
				if isbns[i] {
					Warn("ISBN %s already in library", originalISBN)
					return false
				}
				isbns[i] = true
			}
			for _, i := range book.ISBN13 {
				if isbns[i] {
					Warn("ISBN %s already in library", originalISBN)
					return false
				}
				isbns[i] = true
			}

			lib.Entries = append(lib.Entries, *book)

			addmsg := originalISBN
			if book.Title != "" {
				addmsg += " - " + book.Title
			}
			Print("Added %s", addmsg)
			return true
		}
		return false
	}

	if opts.Positional.ISBN != "" {
		book := getBook(opts.Positional.ISBN, opts)
		addBook(book, opts.Positional.ISBN)
	} else if opts.FilePath != "" || opts.UseStdin {
		var importFile *os.File
		if opts.FilePath != "" {
			Debug("Reading import file %s", opts.FilePath)
			f, err := os.Open(opts.FilePath)
			if err != nil {
				Error("Filed to open %s: %s", opts.FilePath, err.Error())
				return
			}
			defer f.Close()
			importFile = f
		} else {
			Debug("Reading Stdin")
			importFile = os.Stdin
		}

		s := bufio.NewScanner(importFile)
		for s.Scan() {
			isbn := strings.TrimSpace(strings.SplitN(s.Text(), "#", 2)[0])
			if isbn != "" {
				book := getBook(isbn, opts)
				if !addBook(book, isbn) && opts.StopOnError {
					break
				}
			}
		}

		if err := s.Err(); err != nil {
			Error(err.Error())
			return
		}

	} else {
		Error("ISBN or file must be supplied")
		return
	}

	if err := LibrarySave(lib); err != nil {
		Error(err.Error())
		return
	}
}

func getBookByISBN(isbn string, opts *AddOptions) *LibraryEntryV1 {
	// verify ISBN
	isbn = strings.ReplaceAll(isbn, "-", "")

	isbnToTest := isbn
	if len(isbn) == 10 && isbn[9] == 'X' {
		isbnToTest = isbn[0:8]
	}

	if _, err := strconv.Atoi(isbnToTest); err != nil {
		Warn("Invalid ISBN: '%s' - Non-numeric value", isbn)
		return nil
	}

	if len(isbn) == 9 && !opts.NoImpliedISBN10 { // 1967 version => ISO 2108
		Warn("ISBN length of 9 implied leading 0 for ISBN-10")
		isbn = "0" + isbn
	} else if len(isbn) != 10 && len(isbn) != 13 {
		Warn("Invalid ISBN: %s must be length of 10 or 13 characters", isbn)
		return nil
	}

	// OpenLibrary lookup
	ol, err := OLGetByISBN(isbn)
	if err == nil {
		return NewLibraryEntry(ol)
	}

	Warn(err.Error())

	newEntry := LibraryEntryV1{}
	if len(isbn) == 10 {
		newEntry.ISBN10 = []string{isbn}
	} else if len(isbn) == 13 {
		newEntry.ISBN13 = []string{isbn}
	}

	return &newEntry
}

func getBookByOLID(olid string, opts *AddOptions) *LibraryEntryV1 {
	// OpenLibrary lookup
	ol, err := OLGetByOLID(olid)
	if err == nil {
		return NewLibraryEntry(ol)
	}

	Warn(err.Error())

	return &LibraryEntryV1{
		OLID: olid,
	}
}
