package main

import "testing"

type LCCTest struct {
	LCC1   string
	LCC2   string
	Result bool
}

func TestLCCLess(t *testing.T) {
	tests := []LCCTest{
		// Basic
		LCCTest{"M301", "M301", false},
		LCCTest{"M301", "M302", true},
		LCCTest{"M302", "M301", false},

		// Smaller char parts
		LCCTest{"M301", "MX301", true},
		LCCTest{"MX301", "M301", false},

		// numeric tests
		LCCTest{"M04", "M10", true},
		LCCTest{"M4", "M10", true},

		// multi part
		LCCTest{"M301.123", "M301.222", true},
		LCCTest{"M301.222", "M301.123", false},
		LCCTest{"M301 AF", "M301 AA", false},
		LCCTest{"M301 AA", "M301 AF", true},

		LCCTest{"BV5055 .F57 1997", "BV'4639'L45'1991", false},
		LCCTest{"BV'4639'L45'1991", "BV835.I54 2018", false},
		LCCTest{"BV5055 .F57 1997", "BV835.I54 2018", false},

		LCCTest{"BV'4639'L45'1991", "BV5055 .F57 1997", true},
		LCCTest{"BV835.I54 2018", "BV'4639'L45'1991", true},
		LCCTest{"BV835.I54 2018", "BV5055 .F57 1997", true},
	}

	for _, test := range tests {
		// t.Logf("'%s' < '%s' => %t", test.LCC1, test.LCC2, test.Result)
		if lessThan(test.LCC1, test.LCC2) != test.Result {
			t.Errorf("Failed case '%s' < '%s' => %t", test.LCC1, test.LCC2, test.Result)
		}
	}
}

func lessThan(lcc1, lcc2 string) bool {
	lib := LibraryV1{
		Entries: []LibraryEntryV1{
			LibraryEntryV1{
				LCC: lcc1,
			},
			LibraryEntryV1{
				LCC: lcc2,
			},
		},
	}

	lt := lccLessThan(&lib)

	return lt(0, 1)
}
