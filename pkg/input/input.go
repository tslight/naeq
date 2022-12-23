package input

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func Get(args []string, prompt string) (string, error) {
	var scanner *bufio.Scanner
	if len(args) > 0 {
		scanner = bufio.NewScanner(strings.NewReader(strings.Join(args, " ")))
	} else {
		in := os.Stdin
		i, err := in.Stat()
		if err != nil {
			return "", err
		}
		size := i.Size()
		if size == 0 {
			fmt.Print(prompt)
		}
		scanner = bufio.NewScanner(os.Stdin)
	}

	scanner.Scan()
	s := scanner.Text()
	if err := scanner.Err(); err != nil {
		return "", err
	}

	return s, nil
}
