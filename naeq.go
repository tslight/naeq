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
	"unicode"
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

func GetNaeq(s string) int {
	// https://gosamples.dev/remove-non-alphanumeric/
	var nonAlphanumericRegex = regexp.MustCompile(`[^a-zA-Z0-9]+`)
	s = nonAlphanumericRegex.ReplaceAllString(s, "")
	s = strings.ToLower(s)

	value := 0
	var v int
	for _, c := range s {
		if unicode.IsNumber(c) {
			// https://itecnote.com/tecnote/go-convert-rune-to-int/
			// https://stackoverflow.com/a/21322694/11133327 - Bizarre!
			v = int(c - '0')
		} else {
			v = int(c-'a')*19%26 + 1
		}
		value += v
	}

	return value
}

func GetBook(path string) map[string]interface{} {
	var book map[string]interface{}
	if path != "" {
		var err error
		clr.Print(clr.Yel, fmt.Sprintf("Loading %s... ", path))
		book, err = jsn.FromFile(path)
		if err != nil {
			clr.Print(clr.Red, "FAILED! :-(\n")
			log.Fatalln(err)
		}
		clr.Print(clr.Grn, "Done! :-)\n")
	} else {
		clr.Print(clr.Yel, "Loading Liber Al Vegis... ")
		json.Unmarshal([]byte(liberAlBytes), &book)
		clr.Print(clr.Grn, "Done! :-)\n")
	}
	return book
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
	value := GetNaeq(input.Text())
	if err := input.Err(); err != nil {
		log.Fatalln(err)
	}

	book := GetBook(path)

	s := strconv.Itoa(value)
	clr.Printf(clr.Yel, "NAEQ Sum:%s %s\n", clr.Grn, s)

	matches := book[s]
	m := reflect.ValueOf(matches)
	numberOfMatches := m.Len()
	clr.Printf(clr.Yel, "Matches: %s%d\n", clr.Grn, numberOfMatches)
	for i := 0; i < numberOfMatches; i++ {
		if count > 0 && i >= count {
			break
		}
		fmt.Println(m.Index(i))
	}
}
