package main

import (
	_ "embed"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"reflect"
	"strconv"
	"unicode"
)

//go:embed cipher.json
var cipherBytes []byte

//go:embed books/liber-AL.json
var liberAlBytes []byte

const (
	Black  = "\u001b[1;30m"
	Red    = "\u001b[1;31m"
	Green  = "\u001b[1;32m"
	Yellow = "\u001b[1;33m"
	Blue   = "\u001b[1;34m"
	Cyan   = "\u001b[1;36m"
	Reset  = "\u001b[0m"
)

func colorise(color string, msg string) {
	fmt.Print(color, msg, Reset)
}

// https://stackoverflow.com/a/36922225/11133327
func isJson(str string) bool {
	var js json.RawMessage
	return json.Unmarshal([]byte(str), &js) == nil
}

func parseJsonFromFile(path string) (map[string]interface{}, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	defer f.Close()

	byteValue, err := ioutil.ReadAll(f)
	if err != nil {
		return nil, err
	}

	// fmt.Print(string(byteValue))
	if !isJson(string(byteValue)) {
		return nil, fmt.Errorf("%s is not a valid JSON file\n", path)
	}

	var result map[string]interface{}
	json.Unmarshal([]byte(byteValue), &result)

	return result, nil
}

func naeqValueFromString(s string) int {
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

	var book map[string]interface{}
	if path != "" {
		var err error
		colorise(Yellow, fmt.Sprintf("Loading %s... ", path))
		book, err = parseJsonFromFile(path)
		if err != nil {
			colorise(Red, "FAILED! :-(\n")
			fmt.Println(err)
			flag.Usage()
			os.Exit(1)
		}
		colorise(Green, "Done! :-)\n")
	} else {
		colorise(Yellow, "Loading Liber Al Vegis...")
		json.Unmarshal([]byte(liberAlBytes), &book)
		colorise(Green, "Done! :-)\n")
	}

	value := 0
	remainingArgs := flag.NArg()
	if remainingArgs < 1 {
		fmt.Printf("%sNo words to process!%s\n", Red, Reset)
		flag.Usage()
		os.Exit(1)
	}

	for i := 0; i < remainingArgs; i++ {
		arg := flag.Arg(i)
		v := naeqValueFromString(arg)
		// fmt.Printf("%s NAEQ Sum: %d\n", arg, v)
		value += v
	}

	s := strconv.Itoa(value)
	fmt.Printf("%sNAEQ Sum:%s %s\n%s", Yellow, Green, s, Reset)

	matches := book[s]
	m := reflect.ValueOf(matches)
	numberOfMatches := m.Len()
	fmt.Printf("%sMatches:%s %d\n%s", Yellow, Green, numberOfMatches, Reset)
	for i := 0; i < numberOfMatches; i++ {
		if count > 0 && i >= count {
			break
		}
		fmt.Println(m.Index(i))
	}
}
