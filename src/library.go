package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"sort"
	"strconv"
	"unicode"
)

type LibrarySort string

const (
	SortLCC    LibrarySort = "LCC"
	SortAuthor             = "Author"
)

type Library struct {
	V1 *LibraryV1 `json:"v1"`
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

	if lib.V1.SortedBy != SortLCC {
		return nil, fmt.Errorf("Sort type %s not supported by this version of tlib", lib.V1.SortedBy)
	}

	if GlobalOpts.AlwaysPretty {
		lib.V1.AlwaysPretty = true
	} else if GlobalOpts.NotAlwaysPretty {
		lib.V1.AlwaysPretty = false
	}

	return lib.V1, nil
}

func LibraryExists() bool {
	stat, err := os.Stat(LibraryPath())
	return err == nil && !stat.IsDir()
}

func LibraryPath() string {
	return GlobalOpts.FilePath
}

func LibrarySave(lib *LibraryV1) error {
	if lib == nil {
		return errors.New("nil Library")
	}
	Debug("Saving library to %s", LibraryPath())

	func() {
		defer func() {
			if err := recover(); err != nil {
				Error("Error sorting library: ", err)
			}
		}()
		sort.Slice(lib.Entries, lccLessThan(lib))
	}()

	toSave := Library{
		V1: lib,
	}

	var data []byte
	var err error
	if lib.AlwaysPretty || GlobalOpts.Pretty {
		data, err = json.MarshalIndent(toSave, "", "  ")
	} else {
		data, err = json.Marshal(toSave)
	}

	if err != nil {
		return err
	}

	return ioutil.WriteFile(LibraryPath(), data, 0644)
}

func lccLessThan(lib *LibraryV1) func(i, j int) bool {
	return func(i, j int) bool {
		le := lib.Entries[i]
		re := lib.Entries[j]

		llcc := le.LCC
		rlcc := re.LCC
		if llcc == "" {
			llcc = le.SelfClassification
		}
		if rlcc == "" {
			rlcc = re.SelfClassification
		}

		if llcc == rlcc {
			if le.Title == "" {
				return true
			}

			return le.Author < re.Author
		}

		lhs := []rune(llcc)
		rhs := []rune(rlcc)

		if len(rhs) == 0 {
			return false
		}
		if len(lhs) == 0 {
			return true
		}

		li := 0
		ri := 0

		next := func(i int, l []rune) int {
			for ir := i + 1; ir < len(l); ir++ {
				c := l[ir]
				if unicode.IsLetter(c) || unicode.IsNumber(c) {
					return ir
				}
			}
			return -1
		}
		digitLen := func(i int, l []rune) int {
			for dl := 0; i+dl < len(l); dl += 1 {
				if !unicode.IsDigit(l[i+dl]) {
					return dl
				}
			}
			return len(l) - i
		}

		for true {
			// char compare
			for true {
				lc := lhs[li]
				rc := rhs[ri]
				lNum := unicode.IsDigit(lc)
				rNum := unicode.IsDigit(rc)

				if lNum && rNum {
					break // int compare
				} else if lNum {
					return true // 1 < M
				} else if rNum {
					return false // M < 1
				} else if lc < rc {
					return true // A < M
				} else if lc > rc {
					return false // M < A
				}

				// lc == rc
				li = next(li, lhs)
				ri = next(ri, rhs)
				if ri < 0 {
					return false // "AB" < "A" or "A" < "A"
				} else if li < 0 {
					return true // "A" < "AB"
				}
			}

			// int compare
			llen := digitLen(li, lhs)
			rlen := digitLen(ri, rhs)

			ln, _ := strconv.Atoi(string(lhs[li : li+llen]))
			rn, _ := strconv.Atoi(string(rhs[ri : ri+rlen]))

			if ln < rn {
				return true
			} else if rn < ln {
				return false
			}
			li = next(li+llen, lhs)
			ri = next(ri+rlen, rhs)
			if ri < 0 {
				return false // "AB" < "A" or "A" < "A"
			} else if li < 0 {
				return true // "A" < "AB"
			}
		}

		return false
	}
}
