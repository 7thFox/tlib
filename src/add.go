package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type AddOptions struct {
	FilePath        string `short:"f" long:"file" description:"File to load multiple books from"`
	NoImpliedISBN10 bool   `long:"no-implied-isbn10" description:"Disables implied leading 0 for ISBN's with only 9 characters"`
	StopOnError     bool   `short:"s" long:"stop-on-error" description:"Stops file add after first error"`
	UseStdin        bool   `long:"stdin" description:"Read from Stdin"`

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

	addBook := func(book *LibraryEntryV1) bool {
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
				return true
			}
		}
		return false
	}

	if opts.Positional.ISBN != "" {
		book := getBook(opts.Positional.ISBN, opts)
		addBook(book)
	} else if opts.FilePath != "" || opts.UseStdin {
		var f *os.File
		if opts.FilePath != "" {
			f, err := os.Open(opts.FilePath)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Filed to open %s: %s\n", opts.FilePath, err.Error())
			}
			defer f.Close()
		} else {
			f = os.Stdin
		}

		s := bufio.NewScanner(f)
		for s.Scan() {
			isbn := strings.TrimSpace(strings.SplitN(s.Text(), "#", 2)[0])
			if isbn != "" {
				book := getBook(isbn, opts)
				if !addBook(book) && opts.StopOnError {
					break
				}
			}
		}

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
