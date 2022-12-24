package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/tslight/naeq/assets/books"
	"github.com/tslight/naeq/pkg/alw"
	"github.com/tslight/naeq/pkg/efs"
	j "github.com/tslight/naeq/pkg/json"
	"io"
	"log"
	"net/http"
)

type Error struct {
	Code   int    `json:"code"`
	Type   string `json:"domain"`
	Reason string `json:"message"`
}

type Query struct {
	Book  string `json:"book"`
	Words string `json:"words"`
}

func logRequest(r *http.Request) {
	b, err := io.ReadAll(r.Body)
	if err != nil {
		log.Println(err)
	}
	// Replace the body with a new reader after reading from the original
	r.Body = io.NopCloser(bytes.NewBuffer(b))
	ip := r.RemoteAddr
	// ip := r.Header.Get("Client-Ip")
	agent := r.Header.Get("User-Agent")
	log.Printf(
		"%s to %s from %s at %s\n", r.Method, r.URL.Path, agent, ip,
	)
	bstr := string(b)
	if bstr != "" {
		log.Println(bstr)
	}
}

func DecodeRequest(r *http.Request, i interface{}) *Error {
	decoder := json.NewDecoder(r.Body)
	defer r.Body.Close()
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&i); err != nil {
		reason := fmt.Sprintf("Invalid JSON: %v", err)
		return &Error{
			Code:   400,
			Type:   "Bad Request",
			Reason: reason,
		}
	}
	return nil
}

func main() {
	log.Printf("Synchronicity engines started... Awaiting requests...")
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		logRequest(r)
		if r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}
		switch r.Method {
		case http.MethodGet:
			bookNames, err := efs.GetBaseNamesSansExt(&books.EFS)
			if err != nil {
				log.Println(err)
			}
			fmt.Fprintf(w, "Do what thou wilt!\n%v", bookNames)
		case http.MethodPost:
			var query Query
			if err := DecodeRequest(r, &query); err != nil {
				log.Println(err)
				json.NewEncoder(w).Encode(err)
				break
			}
			i, err := alw.GetSum(query.Words)
			if err != nil {
				log.Println(err)
				json.NewEncoder(w).Encode(err)
				break
			}
			log.Printf("NAEQ Sum: %d", i)
			book, err := j.FromEFSPath(books.EFS, query.Book)
			if err != nil {
				log.Println(err)
				json.NewEncoder(w).Encode(err)
				break
			}
			log.Printf("%s = %s", query.Book, book["name"])
			// log.Printf("%#v", query)
			matches := alw.GetMatches(i, book)
			// log.Printf("%#v", matches)
			json.NewEncoder(w).Encode(matches)
			log.Printf("Successfully returned %d matches! :-)", len(matches))
		case http.MethodPut:
			// Update an existing record.
		case http.MethodDelete:
			// Remove the record.
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	log.Fatal(http.ListenAndServe(":10000", nil))
}
