package main

import (
	"bufio"
	"embed"
	// "encoding/json"
	"flag"
	"fmt"
	"github.com/tslight/naeq/clr"
	"github.com/tslight/naeq/jsn"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

//go:embed books/*.json
var books embed.FS

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

func GetBookFromPath(path string) map[string]interface{} {
	var book map[string]interface{}
	var err error
	book, err = jsn.FromPath(path)
	if err != nil {
		log.Fatalln(err)
	}
	return book
}

func GetBookFromEFSPath(efs embed.FS, path string) map[string]interface{} {
	var book map[string]interface{}
	var err error
	book, err = jsn.FromEFSPath(efs, path)
	if err != nil {
		log.Fatalln(err)
	}
	return book
}

func GetMatches(s string, b map[string]interface{}) ([]interface{}, string, error) {
	stats := clr.Sprintf(clr.Yel, "Words:%s %s\n", clr.Grn, s)

	value, err := GetNaeq(s)
	if err != nil {
		return nil, "", err
	}
	stats += clr.Sprintf(clr.Yel, "NAEQ Sum:%s %d\n", clr.Grn, value)

	stats += clr.Sprintf(clr.Yel, "Book:%s %v\n", clr.Grn, b["name"])

	key := strconv.Itoa(value)
	matches := reflect.ValueOf(b[key]).Interface().([]interface{})
	stats += clr.Sprintf(clr.Yel, "Matches: %s%d\n", clr.Grn, len(matches))

	return matches, stats, err
}

func getAllFilenames(efs *embed.FS) (files []string, err error) {
	if err := fs.WalkDir(efs, ".", func(path string, d fs.DirEntry, err error) error {
		if d.IsDir() {
			return nil
		}
		files = append(files, path)
		return nil
	}); err != nil {
		return nil, err
	}
	return files, nil
}

func getBookNames(withLongName bool) string {
	bookNames, err := getAllFilenames(&books)
	var s string
	if err != nil {
		log.Fatalln(err)
	}
	for _, v := range bookNames {
		liberName := strings.TrimSuffix(filepath.Base(v), filepath.Ext(v))
		s += clr.Sprint(clr.Yel, liberName)
		if withLongName {
			book := GetBookFromEFSPath(books, v)
			s += clr.Sprintf(clr.Grn, " (%s)", book["name"])
		}
		s += "\n"
	}
	return s
}

func main() {
	var count int
	var path, efsBook string
	var raw, list, long bool
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s [options...] <words>:\n", os.Args[0])
		flag.PrintDefaults()
		// fmt.Fprintf(os.Stderr, getBookNames(false))
	}
	flag.IntVar(&count, "n", 0, "number of matches to show")
	flag.StringVar(&path, "p", "", "path to alternative book")
	flag.StringVar(&efsBook, "b", "liber-al", "embedded book")
	flag.BoolVar(&raw, "r", false, "display raw unformatted output")
	flag.BoolVar(&list, "l", false, "list embedded books")
	flag.BoolVar(&long, "L", false, "list embedded books with name")
	flag.Parse()

	if list || long {
		fmt.Print(getBookNames(long))
		return
	}

	input, err := GetInput(flag.Args())
	if err != nil {
		log.Fatalln(err)
	}

	input.Scan()
	words := input.Text()
	if err := input.Err(); err != nil {
		log.Fatalln(err)
	}

	var book map[string]interface{}
	if path != "" {
		book = GetBookFromPath(path)
	} else {
		book = GetBookFromEFSPath(books, fmt.Sprint("books/", efsBook, ".json"))
	}
	matches, stats, err := GetMatches(words, book)
	if err != nil {
		log.Fatalln(err)
	}

	if !raw {
		fmt.Print(stats)
	}

	for k, v := range matches {
		if count > 0 && k >= count {
			break
		}
		if raw {
			fmt.Println(v)
		} else {
			fmt.Printf("%d: %s\n", k+1, v)
		}
	}
}
