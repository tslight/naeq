package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/tslight/naeq/assets/books"
	"github.com/tslight/naeq/pkg/alw"
	j "github.com/tslight/naeq/pkg/json"
)

var (
	port    = flag.Int("p", 8080, "Port to listen on")
	version = flag.Bool("v", false, "print version info")
)

var Version = "unknown"

const defaultBook = "liber-al.json"

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
	if r.URL.RawQuery != "" {
		log.Printf("Query Params: %s", r.URL.RawQuery)
	}
	bstr := string(b)
	if bstr != "" {
		log.Println(bstr)
	}
}

func aboutHandler(w http.ResponseWriter, r *http.Request) {
	var scheme string
	if r.TLS == nil {
		scheme = "http"
	} else {
		scheme = "https"
	}
	fmt.Fprintf(w, `
DO WHAT THOU WILT!

The Secret Cipher of the UFOnauts as an API, because ¯\_(ツ)_/¯

https://github.com/tslight/naeq

curl -X GET  %[1]s://%[2]s?words=hellier
curl -X GET  %[1]s://%[2]s?words=hellier&book=liber-i.json
curl -X POST %[1]s://%[2]s -d '{"words": "hellier"}'
curl -X POST %[1]s://%[2]s -d '{"book": "liber-x.json", "words": "hellier"}'
`, scheme, r.Host)
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
		words := r.URL.Query().Get("words")
		if words != "" {
			book = r.URL.Query().Get("book")
			if book == "" {
				book = defaultBook
			}
			response, err = buildResponse(words, book)
			if err != nil {
				log.Println(err)
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		} else {
			aboutHandler(w, r)
			return
		}
	case http.MethodPost:
		data := Data{Book: defaultBook}
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
	log.Println("Synchronicity engines starting...")
	http.HandleFunc("/", Handler)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", *port), nil))
}
