package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

type Article struct {
	Title   string `json:"Title"`
	Desc    string `json:"desc"`
	Content string `json:"content"`
}

var Articles []Article

func logRequest(r *http.Request) {
	ip := r.Header.Get("Client-Ip")
	agent := r.Header.Get("User-Agent")
	log.Printf(
		"%s from %s on %s requested %s\n", r.Method, ip, agent, r.URL.Path,
	)
}

func articles(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(Articles)
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
			fmt.Fprintf(w, "Do what thou wilt!\n")
		case http.MethodPost:
			reqBody, _ := ioutil.ReadAll(r.Body)
			fmt.Fprintf(w, "%+v", string(reqBody))
			log.Print(string(reqBody))
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
