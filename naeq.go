package main

import (
	"bufio"
	_ "embed"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/tslight/naeq/clr"
	"github.com/tslight/naeq/jsn"
	"log"
	"os"
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

//go:embed books/liber-AL.json
var liberAlBytes []byte

func GetInput(args []string) (*bufio.Scanner, error) {
	if len(args) > 0 {
		return bufio.NewScanner(strings.NewReader(strings.Join(args, " "))), nil
	}
	in := os.Stdin
	i, err := in.Stat()
	if err != nil {
		return nil, err
	}
	size := i.Size()
	if size == 0 {
		fmt.Print("Enter text: ")
	}

	return bufio.NewScanner(os.Stdin), nil
}

func SumNumbersInString(s string) (int, error) {
	value := 0
	numericRegex := regexp.MustCompile(`\d+`)
	matches := numericRegex.FindAllString(s, -1)
	for _, m := range matches {
		v, err := strconv.Atoi(m)
		if err != nil {
			return 0, err
		}
		value += v
	}
	return value, nil
}

func GetNaeq(s string) (int, error) {
	// sum all numbers first
	value, err := SumNumbersInString(s)
	if err != nil {
		return 0, err
	}

	nonAlphaRegex := regexp.MustCompile(`[^a-zA-Z]+`)
	s = nonAlphaRegex.ReplaceAllString(s, "")
	s = strings.ToLower(s)

	for _, c := range s {
		value += int(c-'a')*19%26 + 1
	}

	return value, nil
}

func GetBook(path string) map[string]interface{} {
	var book map[string]interface{}
	if path != "" {
		var err error
		book, err = jsn.FromFile(path)
		if err != nil {
			log.Fatalln(err)
		}
	} else {
		json.Unmarshal([]byte(liberAlBytes), &book)
	}
	return book
}

func GetMatches(s string, p string) ([]interface{}, error) {
	value, err := GetNaeq(s)
	if err != nil {
		return nil, err
	}

	clr.Printf(clr.Yel, "NAEQ Sum:%s %d\n", clr.Grn, value)

	book := GetBook(p)
	clr.Printf(clr.Yel, "Book:%s %v\n", clr.Grn, book["name"])

	key := strconv.Itoa(value)

	matches := reflect.ValueOf(book[key]).Interface().([]interface{})
	clr.Printf(clr.Yel, "Matches: %s%d\n", clr.Grn, len(matches))

	return matches, err
}

func main() {
	var count int
	var path string
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s [options...] <words to process>:\n", os.Args[0])
		flag.PrintDefaults()
	}
	flag.IntVar(&count, "n", 0, "number of matches to show")
	flag.StringVar(&path, "p", "", "path to alternative book")
	flag.Parse()
	input, err := GetInput(flag.Args())
	if err != nil {
		log.Fatalln(err)
	}

	input.Scan()
	words := input.Text()
	if err := input.Err(); err != nil {
		log.Fatalln(err)
	}
	clr.Printf(clr.Yel, "Words:%s %s\n", clr.Grn, words)

	matches, err := GetMatches(words, path)
	if err != nil {
		log.Fatalln(err)
	}

	for k, v := range matches {
		if count > 0 && k >= count {
			break
		}
		fmt.Printf("%d: %s\n", k+1, v)
	}
}
