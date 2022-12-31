package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/tslight/naeq/assets/books"
	"github.com/tslight/naeq/pkg/alw"
	j "github.com/tslight/naeq/pkg/json"
	"io"
	"log"
	"net/http"
	"strings"
)

var about = `
DO WHAT THOU WILT!

The Secret Cipher of the UFOnauts as an API, because ¯\_(ツ)_/¯

https://github.com/tslight/naeq
`

var (
	port    = flag.Int("p", 8080, "Port to listen on")
	version = flag.Bool("v", false, "print version info")
)

var Version = "unknown"

type Data struct {
	Book  string `json:"book"`
	Words string `json:"words"`
}

type Response struct {
	Book       interface{}   `json:"book"`
	Sum        int           `json:"sum"`
	MatchCount int           `json:"match_count"`
	Matches    []interface{} `json:"matches"`
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

func buildResponse(words string, book string) (interface{}, error) {
	i, err := alw.GetSum(words)
	if err != nil {
		return nil, err
	}
	log.Printf("NAEQ Sum: %d", i)
	b, err := j.FromEFSPath(books.EFS, book)
	if err != nil {
		return nil, err
	}
	log.Printf("%s = %s", book, b["name"])
	matches := alw.GetMatches(i, b)
	log.Printf("Successfully found %d matches! :-)", len(matches))
	response := Response{
		Book:       b["name"],
		Sum:        i,
		MatchCount: len(matches),
		Matches:    matches,
	}
	return response, nil
}

func Handler(w http.ResponseWriter, r *http.Request) {
	var response interface{}
	var err error
	var book string

	logRequest(r)
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	switch r.Method {
	case http.MethodGet:
		query := r.URL.Query()
		if len(query["words"]) > 0 {
			words := strings.Join(query["words"], " ")
			if len(query["book"]) > 0 {
				book = strings.Join(query["book"], " ")
			} else {
				book = "liber-al.json"
			}
			response, err = buildResponse(words, book)
			if err != nil {
				log.Println(err)
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		} else {
			fmt.Fprint(w, about)
			return
		}
	case http.MethodPost:
		var data Data
		decoder := json.NewDecoder(r.Body)
		defer r.Body.Close()
		decoder.DisallowUnknownFields()
		if err := decoder.Decode(&data); err != nil {
			log.Println(err)
			http.Error(
				w, fmt.Sprintf("Bad Request: %s", err.Error()), http.StatusBadRequest,
			)
			return
		}
		response, err = buildResponse(data.Words, data.Book)
		if err != nil {
			log.Println(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		log.Println(err)
	}
}

func main() {
	flag.Parse()
	if *version {
		fmt.Println(Version)
		return
	}
	log.Printf("Synchronicity engines started...")
	http.HandleFunc("/", Handler)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", *port), nil))
}
