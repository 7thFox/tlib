package main

type LibraryV1 struct {
	AlwaysPretty bool             `json:"pretty"`
	SortedBy     LibrarySort      `json:"sorted_by"`
	Entries      []LibraryEntryV1 `json:"books"`
}

type LibraryEntryV1 struct {
	OLID   string   `json:"ol_id,omitempty"`
	ISBN10 []string `json:"isbn_10,omitempty"`
	ISBN13 []string `json:"isbn_13,omitempty"`
	LCC    string   `json:"lcc,omitempty"`
	Title  string   `json:"title,omitempty"`
	Author string   `json:"author,omitempty"`

	SelfClassification string `json:"self_lcc,omitempty"`
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
	e := LibraryEntryV1{}
	e.Title = ol.Title
	e.Author = ol.ByStatement
	if e.Author == "" {
		for i, a := range ol.Authors {
			if i == len(ol.Authors)-1 {
				e.Author += a.Name
			} else if i == len(ol.Authors)-2 {
				e.Author += a.Name + ", and "
			} else {

				e.Author += a.Name + ", "
			}
		}
	}

	if len(ol.Classifications.LCClassifications) > 0 {
		e.LCC = ol.Classifications.LCClassifications[0]
	} else {
		Warn("No LC Classification found on OpenLibrary")
	}
	if len(ol.Identifiers.OpenLibrary) > 0 {
		e.OLID = ol.Identifiers.OpenLibrary[0]
	} else {
		Warn("No OpenLibrary ID found on OpenLibrary")
	}

	e.ISBN10 = ol.Identifiers.ISBN_10
	e.ISBN13 = ol.Identifiers.ISBN_13

	return &e
}

func (e *LibraryEntryV1) Update(ol *OpenLibraryBook) {
	if len(ol.Identifiers.OpenLibrary) > 0 {
		e.OLID = ol.Identifiers.OpenLibrary[0]
	}

	if ol.Title != "" {
		e.Title = ol.Title
	}
	if ol.ByStatement != "" {
		e.Author = ol.ByStatement
	}
	if e.Author == "" {
		for i, a := range ol.Authors {
			if i == len(ol.Authors)-1 {
				e.Author += a.Name
			} else if i == len(ol.Authors)-2 {
				e.Author += a.Name + ", and "
			} else {

				e.Author += a.Name + ", "
			}
		}
	}
	if len(ol.Classifications.LCClassifications) > 0 {
		e.LCC = ol.Classifications.LCClassifications[0]
	}

	if e.LCC != "" {
		e.SelfClassification = ""
	}

	e.ISBN10 = ol.Identifiers.ISBN_10
	e.ISBN13 = ol.Identifiers.ISBN_13
}
