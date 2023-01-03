package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"net/http"
	"os"

	"github.com/tslight/naeq/assets/books"
	"github.com/tslight/naeq/pkg/alw"
	"github.com/tslight/naeq/pkg/efs"
	j "github.com/tslight/naeq/pkg/json"
	"github.com/tslight/naeq/pkg/log"
	"time"
)

var about = `DO WHAT THOU WILT!

The Secret Cipher of the UFOnauts as an API, because ¯\_(ツ)_/¯

https://github.com/tslight/naeq
`

var (
	port    = flag.String("p", "8080", "Port to listen on")
	version = flag.Bool("v", false, "print version info")
)

// Version set at compile time from envvar with -X main.Version=$(VERSION)
var Version = "unknown"

const defaultBook = "liber-al"

type Data struct {
	Book  string `json:"book"`
	Words string `json:"words"`
}

type Response struct {
	Liber      interface{}   `json:"liber"`
	Book       interface{}   `json:"book"`
	Sum        int           `json:"sum"`
	MatchCount int           `json:"match_count"`
	Matches    []interface{} `json:"matches"`
}

func aboutHandler(w http.ResponseWriter, r *http.Request) error {
	var scheme string
	if r.TLS == nil {
		scheme = "http"
	} else {
		scheme = "https"
	}

	bookNames, err := efs.GetBaseNamesSansExt(&books.EFS)
	if err != nil {
		return err
	}
	bookStr := ""
	for _, v := range bookNames {
		bookStr += fmt.Sprintln(v)
	}

	fmt.Fprintf(w, about+`
Examples:

curl -X GET  %[1]s://%[2]s?words=hellier
curl -X GET  %[1]s://%[2]s?words=hellier&book=liber-i
curl -X POST %[1]s://%[2]s -d '{"words": "hellier"}'
curl -X POST %[1]s://%[2]s -d '{"book": "liber-x", "words": "hellier"}'

Available Books:

%s
`, scheme, r.Host, bookStr)

	return nil
}

func buildResponse(words string, book string) (interface{}, error) {
	i, err := alw.GetSum(words)
	if err != nil {
		return nil, err
	}
	log.Info.Print("NAEQ: ", i)

	b, err := j.FromEFSPath(books.EFS, fmt.Sprint(book, ".json"))
	if err != nil {
		return nil, err
	}
	log.Info.Printf("BOOK: %s (%s)", b["liber"], b["name"])

	matches := alw.GetMatches(i, b)
	log.Info.Print("MATCHES: ", len(matches))

	response := Response{
		Liber:      b["liber"],
		Book:       b["name"],
		Sum:        i,
		MatchCount: len(matches),
		Matches:    matches,
	}

	log.Debug.Printf("RESPONSE:\n%#v", response)
	return response, nil
}

func handler(w http.ResponseWriter, r *http.Request) {
	var response interface{}
	var err error
	var book string

	log.Request(r)
	if r.URL.Path != "/" {
		log.Error.Print(http.StatusNotFound, " NOT FOUND")
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
				errMsg := fmt.Sprintf(
					"%d INTERNAL SERVER ERROR %v", http.StatusInternalServerError, err,
				)
				log.Error.Print(errMsg)
				http.Error(w, errMsg, http.StatusInternalServerError)
				return
			}
		} else {
			err := aboutHandler(w, r)
			if err != nil {
				log.Error.Println(err)
			}
			return
		}
	case http.MethodPost:
		data := Data{Book: defaultBook}
		decoder := json.NewDecoder(r.Body)
		defer r.Body.Close()
		decoder.DisallowUnknownFields()
		if err := decoder.Decode(&data); err != nil {
			errMsg := fmt.Sprintf("%d BAD REQUEST %v", http.StatusBadRequest, err)
			log.Error.Print(errMsg)
			http.Error(w, errMsg, http.StatusBadRequest)
			return
		}
		response, err = buildResponse(data.Words, data.Book)
		if err != nil {
			errMsg := fmt.Sprintf(
				"%d INTERNAL SERVER ERROR %v", http.StatusInternalServerError, err,
			)
			log.Error.Print(errMsg)
			http.Error(w, errMsg, http.StatusInternalServerError)
			return
		}
	default:
		errMsg := fmt.Sprint(http.StatusMethodNotAllowed, " METHOD NOT ALLOWED")
		log.Error.Print(errMsg)
		http.Error(w, errMsg, http.StatusMethodNotAllowed)
		return
	}

	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		log.Error.Print(err)
	}
}

func main() {
	flag.Parse()
	if *version {
		fmt.Println(Version)
		return
	}

	envPort, envPortPresent := os.LookupEnv("PORT")
	if envPortPresent {
		*port = envPort
	}

	log.Info.Print("Synchronicity engines starting on PORT:", *port)
	http.HandleFunc("/", handler)

	server := &http.Server{
		Addr:              ":" + *port,
		ReadHeaderTimeout: 3 * time.Second,
	}

	err := server.ListenAndServe()
	if err != nil {
		log.Error.Fatal(err)
	}
}
