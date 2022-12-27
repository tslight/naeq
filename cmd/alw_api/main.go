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

type Query struct {
	Book  string `json:"book"`
	Words string `json:"words"`
}

type Response struct {
	Book       interface{}   `json:"book"`
	Sum        int           `json:"sum"`
	MatchCount int           `json:"match_count"`
	Matches    []interface{} `json:matches`
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

func Handler(w http.ResponseWriter, r *http.Request) {
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
		json.NewEncoder(w).Encode(bookNames)
	case http.MethodPost:
		var query Query
		decoder := json.NewDecoder(r.Body)
		defer r.Body.Close()
		decoder.DisallowUnknownFields()
		if err := decoder.Decode(&query); err != nil {
			log.Println(err)
			http.Error(
				w, fmt.Sprintf("Bad Request: %s", err.Error()), http.StatusBadRequest,
			)
			return
		}
		i, err := alw.GetSum(query.Words)
		if err != nil {
			log.Println(err)
			http.Error(
				w,
				fmt.Sprintf("Internal Server Error: %v", err.Error()),
				http.StatusInternalServerError,
			)
			return
		}
		log.Printf("NAEQ Sum: %d", i)
		book, err := j.FromEFSPath(books.EFS, query.Book)
		if err != nil {
			log.Println(err)
			http.Error(
				w,
				fmt.Sprintf("Internal Server Error: %v", err.Error()),
				http.StatusInternalServerError,
			)
			return
		}
		log.Printf("%s = %s", query.Book, book["name"])
		// log.Printf("%#v", query)
		matches := alw.GetMatches(i, book)
		response := Response{
			Book:       book["name"],
			Sum:        i,
			MatchCount: len(matches),
			Matches:    matches,
		}
		// log.Printf("%#v", response)
		json.NewEncoder(w).Encode(response)
		log.Printf("Successfully returned %d matches! :-)", len(matches))
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func main() {
	log.Printf("Synchronicity engines started...")
	http.HandleFunc("/", Handler)
	log.Fatal(http.ListenAndServe(":10000", nil))
}
