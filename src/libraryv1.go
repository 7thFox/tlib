package main

type LibraryV1 struct {
	SortedBy LibrarySort      `json:"sorted_by"`
	Entries  []LibraryEntryV1 `json:"books"`
}

type LibraryEntryV1 struct {
	ISBN10 []string `json:"isbn_10,omitempty"`
	ISBN13 []string `json:"isbn_13,omitempty"`
	LCC    string   `json:"lcc,omitempty"`
	Title  string   `json:"title,omitempty"`
	Author string   `json:"author,omitempty"`

	OpenLibraryURL string `json:"ol_url,omitempty"`
}

func (e LibraryEntryV1) FirstISBN() string {
	for _, i := range e.ISBN13 {
		return i
	}

	for _, i := range e.ISBN10 {
		return i
	}
	return ""
}

func NewLibraryEntry(ol *OpenLibraryBook) *LibraryEntryV1 {
	newEntry := LibraryEntryV1{}
	newEntry.Title = ol.Title
	newEntry.Author = ol.ByStatement
	if len(ol.Classifications.LCClassifications) > 0 {
		newEntry.LCC = ol.Classifications.LCClassifications[0]
	} else {
		Warn("No LC Classification found on OpenLibrary")
	}

	newEntry.ISBN10 = ol.Identifiers.ISBN_10
	newEntry.ISBN13 = ol.Identifiers.ISBN_13

	return &newEntry
}

func (e *LibraryEntryV1) Update(ol *OpenLibraryBook) {
	if ol.OpenLibraryKey != "" {
		e.OpenLibraryURL = "https://openlibrary.org" + ol.OpenLibraryKey
	}

	if ol.Title != "" {
		e.Title = ol.Title
	}
	if ol.ByStatement != "" {
		e.Author = ol.ByStatement
	}
	if len(ol.Classifications.LCClassifications) > 0 {
		e.LCC = ol.Classifications.LCClassifications[0]
	}

	e.ISBN10 = ol.Identifiers.ISBN_10
	e.ISBN13 = ol.Identifiers.ISBN_13
}
