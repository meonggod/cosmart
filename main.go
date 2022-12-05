package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

type Books struct {
	//Key         string  `json:"key"`
	//Name        string  `json:"name"`
	//SubjectType string  `json:"subject_type"`
	//WorkCount   int     `json:"work_count"`
	Works []Works `json:"works"`
}

type Works struct {
	Key   string `json:"key"`
	Title string `json:"title"`
	//EditionCount      int      `json:"edition_count"`
	//CoverId           int      `json:"cover_id"`
	//CoverEditionKey   string   `json:"cover_edition_key"`
	Subject []string `json:"subject"`
	//IaCollection      []string `json:"ia_collection"`
	//Lendinglibrary    bool     `json:"lendinglibrary"`
	//Printdisabled     bool     `json:"printdisabled"`
	//LendingEdition    string   `json:"lending_edition"`
	//LendingIdentifier string   `json:"lending_identifier"`
	//Authors           []struct {
	//	Key  string `json:"key"`
	//	Name string `json:"name"`
	//} `json:"authors"`
	//FirstPublishYear int    `json:"first_publish_year"`
	//Ia               string `json:"ia"`
	//PublicScan       bool   `json:"public_scan"`
	//HasFulltext      bool   `json:"has_fulltext"`
	Availability Availability `json:"availability"`
}

type Availability struct {
	//Status              string      `json:"status"`
	//AvailableToBrowse   bool        `json:"available_to_browse"`
	AvailableToBorrow bool   `json:"available_to_borrow"`
	BorrowedAt        string `json:"borrowed_at"`
	//AvailableToWaitlist bool        `json:"available_to_waitlist"`
	//IsPrintdisabled     bool        `json:"is_printdisabled"`
	//IsReadable          bool        `json:"is_readable"`
	//IsLendable          bool        `json:"is_lendable"`
	//IsPreviewable       bool        `json:"is_previewable"`
	//Identifier          string      `json:"identifier"`
	//Isbn                *string     `json:"isbn"`
	//Oclc                interface{} `json:"oclc"`
	//OpenlibraryWork     string      `json:"openlibrary_work"`
	//OpenlibraryEdition  string      `json:"openlibrary_edition"`
	//LastLoanDate        interface{} `json:"last_loan_date"`
	//NumWaitlist         interface{} `json:"num_waitlist"`
	//LastWaitlistDate    interface{} `json:"last_waitlist_date"`
	//IsRestricted        bool        `json:"is_restricted"`
	//IsBrowseable        bool        `json:"is_browseable"`
	//Src                 string      `json:"__src__"`
}

var books Books

func main() {
	err := books.readJsonFile()
	if err != nil {
		log.Fatalln("Error read json file", err)
	}

	NewRoute()
}

func getBooks() *Books {
	return &Books{}
}

func (books *Books) readJsonFile() error {

	jsonFile, err := os.Open("data.json")
	if err != nil {
		return err
	}
	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)

	json.Unmarshal(byteValue, books)

	return nil
}

func method(next http.Handler, method string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		type response struct {
			Message string
		}

		if r.Method != method {
			json.NewEncoder(w).Encode(response{Message: "Wrong Method"})
			return
		}
		next.ServeHTTP(w, r)

	})
}

func NewRoute() {
	r := http.NewServeMux()

	r.Handle("/book/", method(http.HandlerFunc(bookSchedule), "GET"))
	r.Handle("/list", method(http.HandlerFunc(list), "POST"))

	log.Println("Listening to localhost:8080")
	http.ListenAndServe("127.0.0.1:8080", r)
}

func list(w http.ResponseWriter, r *http.Request) {
	var data []Works
	w.Header().Set("Content-Type", "application/json")

	subject := r.FormValue("subject")

	if len(subject) == 0 {
		json.NewEncoder(w).Encode(books.Works)
		return
	}

	for _, row := range books.Works {
		if getIndex(row.Subject, subject) != -1 {
			data = append(data, Works{
				Key:   row.Key,
				Title: row.Title,
				Availability: Availability{
					AvailableToBorrow: row.Availability.AvailableToBorrow,
					BorrowedAt:        row.Availability.BorrowedAt,
				},
			})
		}
	}

	json.NewEncoder(w).Encode(data)
}

func bookSchedule(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	type response struct {
		Message string
	}
	url := strings.Split(strings.TrimPrefix(fmt.Sprintf(`%v`, r.URL), "/"), "/")
	if len(url) > 2 {
		json.NewEncoder(w).Encode(response{Message: "Internal Server Error"})
		return
	}

	location, _ := time.LoadLocation("Asia/Jakarta")
	timeBook := time.Now().In(location).Format("2006-01-02 15:04:05")

	for i, row := range books.Works {
		if row.Key == url[1] {
			if books.Works[i].Availability.AvailableToBorrow == false {
				json.NewEncoder(w).Encode(response{Message: "Book not available to book"})
				return
			}
			books.Works[i].Availability.AvailableToBorrow = false
			books.Works[i].Availability.BorrowedAt = timeBook
			json.NewEncoder(w).Encode(response{Message: fmt.Sprintf("Book Success at: %v", timeBook)})
			return
		}
	}

	json.NewEncoder(w).Encode(response{Message: "Book not found"})
}

func getIndex(haystack []string, needle string) int {
	for i, row := range haystack {
		if strings.Contains(strings.ToLower(row), strings.ToLower(needle)) {
			return i
		}
	}
	return -1
}
