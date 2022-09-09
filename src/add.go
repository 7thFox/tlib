package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

type AddOptions struct {
	FilePath        string `short:"f" long:"file" description:"File to load multiple books from"`
	NoImpliedISBN10 bool   `long:"no-implied-isbn10" description:"Disables implied leading 0 for ISBN's with only 9 characters"`

	Positional struct {
		ISBN string
	} `positional-args:"true"`
}

func RunAdd(opts *AddOptions) {
	lib, err := LibraryLoad()
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		return
	}

	isbns := make(map[string]bool, len(lib.Entries))
	for _, e := range lib.Entries {
		if e.ISBN != "" {
			isbns[e.ISBN] = true
		}
	}

	addBook := func(book *LibraryEntryV1) {
		if book != nil {
			if isbns[book.ISBN] {
				fmt.Fprintln(os.Stderr, "ISBN already in library")
			} else {
				if book.ISBN != "" {
					isbns[book.ISBN] = true
				}

				lib.Entries = append(lib.Entries, *book)

				addmsg := book.ISBN
				if book.Title != "" {
					addmsg += " - " + book.Title
				}
				fmt.Printf("Added %s\n", addmsg)
			}
		}
	}

	if opts.Positional.ISBN != "" {
		book := getBook(opts.Positional.ISBN, opts)
		addBook(book)
	} else if opts.FilePath == "" {
	} else {
		fmt.Fprintln(os.Stderr, "ISBN or file must be supplied")
		return
	}

	if err := LibrarySave(lib); err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		return
	}
}

func getBook(isbn string, opts *AddOptions) *LibraryEntryV1 {
	// verify ISBN
	isbn = strings.ReplaceAll(isbn, "-", "")
	if _, err := strconv.Atoi(isbn); err != nil {
		fmt.Fprintln(os.Stderr, "Invalid ISBN")
		return nil
	}

	if len(isbn) == 9 && !opts.NoImpliedISBN10 {
		fmt.Fprintln(os.Stderr, "ISBN length of 9 implied leading 0 for ISBN-10")
		isbn = "0" + isbn
	} else if len(isbn) != 10 && len(isbn) != 13 {
		fmt.Fprintln(os.Stderr, "ISBN must be length of 10 or 13 characters")
		return nil
	}

	// OpenLibrary lookup
	newEntry := LibraryEntryV1{
		ISBN: isbn,
	}

	ol, err := OLGetByISBN(newEntry.ISBN)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Failed to find in OpenLibrary: "+err.Error())
		// continue
	} else {
		newEntry.Title = ol.Title
		newEntry.Author = ol.ByStatement
		if len(ol.Classifications.LCClassifications) > 0 {
			newEntry.LCC = ol.Classifications.LCClassifications[0]
		} else {
			fmt.Fprintln(os.Stderr, "No LC Classification found on OpenLibrary")
		}
	}

	return &newEntry
}
