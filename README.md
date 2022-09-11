# tlib - Terminal Library
tlib is a light-weight CLI-based library manager utilizing the OpenLibrary API (https://openlibrary.org/)

# features
- Mass-add lists of ISBNs through file or stdin
- No need to manually add information; everything is imported from OpenLibrary on add
- Easily editable JSON format for book listing
- List books in shelf order by LCC (Library of Congress Classifiction)
- List books with missing information not available on OpenLibrary
- Self-Classification field `self_lcc` for when there's not official classification given (currently only available by editing JSON)
- Add by OpenLibrary ID for books without an ISBN

# todo
- Remove book
- Allow ISBN-10s to end with -X
- Validate ISBN checksums

# example usage

```
$ tlib init
$ tlib add 9780141395876
Added 9780141395876 - Prince
$ tlib add 9780521755511
Added 9780521755511 - The Elements of New Testament Greek
$ tlib add 0441013597
Added 0441013597 - Dune
$ tlib add 9780226738925
Added 9780226738925 - The concept of the political
$ tlib shelf --list
JA74 .S313 2007      - The concept of the political
JC143.M38 2014       - Prince
PA817.W4 2005        - The Elements of New Testament Greek
PS3558.E63 D8 2005   - Dune
```

# man page

As of v1.0

```
NAME
       tlib - Terminal Library

SYNOPSIS
       tlib [OPTIONS]

DESCRIPTION
       Light-weight CLI-based library manager utilizing the OpenLibrary API (https://openlibrary.org/)

OPTIONS
   Application Options
       --lib <default: "tlib.json">
              Library file path

       -p, --pretty
              Write library JSON in human-readable format

       --debug
              Print debug info
              

       -q, --quiet
              Only print error messages

       -v, --version
              Print tlib version

   Help Options
       -h, --help
              Show this help message

COMMANDS
   add
       Add books to library

       Usage: tlib [OPTIONS] add [add-OPTIONS]

          Add books to library

       -f, --file
              File to load multiple books from

       --no-implied-isbn10
              Disables implied leading 0 for ISBN's with only 9 characters

       -s, --stop-on-error
              Stops file add after first error

       --stdin
              Read from Stdin

   Help Options
       -h, --help
              Show this help message

   init
       Initialize a new library

   Help Options
       -h, --help
              Show this help message

   shelf
       Print books and other library-wide operations

       Usage: tlib [OPTIONS] shelf [shelf-OPTIONS]

          Print books and other library-wide operations

       -l, --list

       -a, --all
              Print all books including those with LCC or Title

       --no-info
              Only print books without information (ISBN only)

       --no-class
              Only print books without a LCC

       --resort
              Resort your library

       -f, --find-missing
              Search OpenLibrary for books missing a Title/LCC

       --find-all
              Search OpenLibrary for all books and update

   Help Options
       -h, --help
              Show this help message

```
