package main

import (
	"flag"
	"fmt"
	"github.com/tslight/naeq/assets/books"
	"github.com/tslight/naeq/pkg/alw"
	"github.com/tslight/naeq/pkg/efs"
	"github.com/tslight/naeq/pkg/input"
	"github.com/tslight/naeq/pkg/json"
	"log"
	"os"
)

var (
	count   = flag.Int("n", 0, "number of matches to show")
	path    = flag.String("p", "", "path to alternative book")
	efsBook = flag.String("b", "liber-al", "embedded book")
	list    = flag.Bool("l", false, "list embedded books")
	sum     = flag.Bool("s", false, "only return naeq sum")
	version = flag.Bool("v", false, "print version info")
)

var Version = "unknown"

func usage() {
	fmt.Fprintf(os.Stderr, "Usage: %s [options...] <words>:\n", os.Args[0])
	flag.PrintDefaults()
}

func main() {
	flag.Usage = usage
	flag.Parse()

	if *version {
		fmt.Printf("%s %s\n", os.Args[0], Version)
		return
	}

	if *list {
		bookNames, err := efs.GetBaseNamesSansExt(&books.EFS)
		if err != nil {
			log.Fatalln(err)
		}
		for _, v := range bookNames {
			fmt.Println(v)
		}
		return
	}

	words, err := input.Get(flag.Args(), "Enter words: ")
	if err != nil {
		log.Fatalln(err)
	}

	i, err := alw.GetSum(words)
	if err != nil {
		log.Fatalln(err)
	}

	if *sum {
		fmt.Println(i)
		return
	}

	var book map[string]interface{}
	if *path != "" {
		book, err = json.FromPath(*path)
		if err != nil {
			log.Fatalln(err)
		}
	} else {
		book, err = json.FromEFSPath(books.EFS, fmt.Sprint(*efsBook, ".json"))
		if err != nil {
			log.Fatalln(err)
		}
	}

	matches := alw.GetMatches(i, book)

	for k, v := range matches {
		if *count > 0 && k >= *count {
			break
		}
		fmt.Println(v)
	}
}
