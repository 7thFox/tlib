package main

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
)

func OLGetByISBN(isbn string) (*OpenLibraryBook, error) {
	resp, err := http.Get("https://openlibrary.org/api/books?bibkeys=ISBN:" + isbn + "&jscmd=data&format=json")
	defer resp.Body.Close()
	if err != nil {
		return nil, err
	}
	if resp.StatusCode < 200 || resp.StatusCode > 200 {
		return nil, errors.New("OpenLibrary returned failure response")
	}

	bytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var value map[string]OpenLibraryBook
	err = json.Unmarshal(bytes, &value)
	if err != nil {
		return nil, err
	}

	for _, v := range value {
		return &v, nil
	}
	return nil, errors.New("Not Found on OpenLibrary")
}

type OpenLibraryBook struct {
	Url           string `json:"url"`
	Title         string `json:"title"`
	Subtitle      string `json:"subtitle"`
	ByStatement   string `json:"by_statement"`
	NumberOfPages uint   `json:"number_of_pages"`
	Weight        string `json:"weight"`
	PublishDate   string `json:"publish_date"`

	Authors []struct {
		Name string `json:"name"`
		Url  string `json:"url"`
	} `json:"authors"`

	Identifiers struct {
		ISBN_10   []string `json:"isbn_10"`
		ISBN_13   []string `json:"isbn_13"`
		LCCN      []string `json:"lccn"`
		OCLC      []string `json:"oclc"`
		Goodreads []string `json:"goodreads"`
	} `json:"identifiers"`

	Classifications struct {
		LCClassifications []string `json:"lc_classifications"`
		DeweyDecimalClass []string `json:"dewey_decimal_class"`
	} `json:"classifications"`

	// Subjects []struct {
	// 	Name string `json:"name"`
	// 	Url  string `json:"url"`
	// } `json:"subjects"`

	// SubjectPlaces []struct {
	// 	Name string `json:"name"`
	// 	Url  string `json:"url"`
	// } `json:"subject_places"`

	// SubjectPeople []struct {
	// 	Name string `json:"name"`
	// 	Url  string `json:"url"`
	// } `json:"subject_people"`

	// SubjectTimes []struct {
	// 	Name string `json:"name"`
	// 	Url  string `json:"url"`
	// } `json:"subject_times"`

	// Publishers []struct {
	// 	Name string `json:"name"`
	// } `json:"publishers"`

	// PublishPlaces []struct {
	// 	Name string `json:"name"`
	// } `json:"publish_places"`

	// Excerpts []struct {
	// 	Comment string `json:"comment"`
	// 	Text    string `json:"text"`
	// } `json:"excerpts"`

	// Links []struct {
	// 	Url   string `json:"url"`
	// 	Title string `json:"title"`
	// } `json:"links"`

	// Cover struct {
	// 	Small  string `json:"small"`
	// 	Medium string `json:"medium"`
	// 	Large  string `json:"large"`
	// } `json:"cover"`

	// EBooks struct {
	// 	PreviewUrl string `json:"preview_url"`
	// } `json:"ebooks"`
}
