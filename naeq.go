package main

import (
	_ "embed"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/tslight/naeq/clr"
	"github.com/tslight/naeq/jsn"
	"os"
	"reflect"
	"strconv"
	"unicode"
)

//go:embed cipher.json
var cipherBytes []byte

//go:embed books/liber-AL.json
var liberAlBytes []byte

func NaeqFromString(s string) int {
	var cipher map[string]interface{}
	json.Unmarshal([]byte(cipherBytes), &cipher)
	value := 0
	var v int
	for _, char := range s {
		if unicode.IsNumber(char) {
			// fmt.Printf("%c is a number\n", char)
			// https://itecnote.com/tecnote/go-convert-rune-to-int/
			// https://stackoverflow.com/a/21322694/11133327 - Bizarre!
			v = int(char - '0')
		} else {
			// fmt.Printf("%c is NOT a number\n", char)
			v = int(cipher[string(char)].(float64))
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
			fmt.Println(err)
			flag.Usage()
			os.Exit(1)
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

	value := 0
	remainingArgs := flag.NArg()
	if remainingArgs < 1 {
		clr.Print(clr.Red, "No words to process!\n")
		flag.Usage()
		os.Exit(1)
	}

	book := GetBook(path)

	for i := 0; i < remainingArgs; i++ {
		arg := flag.Arg(i)
		v := NaeqFromString(arg)
		// fmt.Printf("%s NAEQ Sum: %d\n", arg, v)
		value += v
	}

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
