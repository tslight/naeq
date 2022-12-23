package main

import (
	"embed"
	"flag"
	"fmt"
	"github.com/tslight/naeq/pkg/efs"
	"github.com/tslight/naeq/pkg/input"
	"github.com/tslight/naeq/pkg/json"
	"github.com/tslight/naeq/pkg/naeq"
	"log"
	"os"
)

//go:embed assets/*.json
var books embed.FS

func main() {
	var count int
	var path, efsBook string
	var list, sum bool
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s [options...] <words>:\n", os.Args[0])
		flag.PrintDefaults()
	}
	flag.IntVar(&count, "n", 0, "number of matches to show")
	flag.StringVar(&path, "p", "", "path to alternative book")
	flag.StringVar(&efsBook, "b", "liber-al", "embedded book")
	flag.BoolVar(&list, "l", false, "list embedded books")
	flag.BoolVar(&sum, "s", false, "only return naeq sum")
	flag.Parse()

	if list {
		fmt.Print(efs.GetBaseNamesSansExt(&books))
		return
	}

	words, err := input.Get(flag.Args(), "Enter words: ")
	if err != nil {
		log.Fatalln(err)
	}

	i, err := naeq.GetSum(words)
	if err != nil {
		log.Fatalln(err)
	}

	if sum {
		fmt.Println(i)
		return
	}

	var book map[string]interface{}
	if path != "" {
		book, err = json.FromPath(path)
		if err != nil {
			log.Fatalln(err)
		}
	} else {
		book, err = json.FromEFSPath(books, fmt.Sprint("assets/", efsBook, ".json"))
		if err != nil {
			log.Fatalln(err)
		}
	}

	matches := naeq.GetMatches(i, book)

	for k, v := range matches {
		if count > 0 && k >= count {
			break
		}
		fmt.Println(v)
	}
}
