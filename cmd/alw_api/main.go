package main

import (
	"encoding/json"
	"fmt"
	"github.com/tslight/naeq/assets/books"
	"github.com/tslight/naeq/pkg/alw"
	"github.com/tslight/naeq/pkg/efs"
	j "github.com/tslight/naeq/pkg/json"
	"io/ioutil"
	"log"
	"net/http"
)

type Query struct {
	Book  string `json:"book"`
	Words string `json:"words"`
}

func logRequest(r *http.Request) {
	ip := r.Header.Get("Client-Ip")
	agent := r.Header.Get("User-Agent")
	log.Printf(
		"%s from %s on %s requested %s\n", r.Method, ip, agent, r.URL.Path,
	)
}

func main() {
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
			reqBody, err := ioutil.ReadAll(r.Body)
			if err != nil {
				log.Println(err)
			}
			log.Print(string(reqBody))
			var query Query
			json.Unmarshal(reqBody, &query)
			i, err := alw.GetSum(query.Words)
			if err != nil {
				log.Println(err)
			}
			book, err := j.FromEFSPath(books.EFS, query.Book)
			if err != nil {
				log.Println(err)
			}
			matches := alw.GetMatches(i, book)
			json.NewEncoder(w).Encode(matches)
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
