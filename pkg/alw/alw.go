package alw

import (
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

func sumNumbersInString(s string) (int, error) {
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

func GetSum(s string) (int, error) {
	// sum all numbers first
	value, err := sumNumbersInString(s)
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

func GetMatches(sum int, b map[string]any) []any {
	key := strconv.Itoa(sum)
	matches := reflect.ValueOf(b[key]).Interface().([]any)

	return matches
}
