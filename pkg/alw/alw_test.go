package alw

import (
	"fmt"
	"testing"

	"github.com/tslight/naeq/pkg/json"
)

func TestSumNumbersInString(t *testing.T) {
	s := "93 H3lli3r 93"
	got, err := sumNumbersInString(s)
	want := 192
	if err != nil {
		t.Fatalf(`SumNumbersInString(%s) returned %d, instead of %d`, s, err, want)
	}
	if got != want {
		t.Fatalf(`SumNumbersInString(%s) returned %d, instead of %d`, s, got, want)
	}
}

func TestGetSum(t *testing.T) {
	s := "#! Hellier  !#"
	got, err := GetSum(s)
	want := 93
	if err != nil {
		t.Fatalf(`GetSum(%s) returned %d, instead of %d`, s, err, want)
	}
	if got != want {
		t.Fatalf(`GetSum(%s) returned %d, instead of %d`, s, got, want)
	}
}

func TestGetSumWithNumbers(t *testing.T) {
	s := "31 #! H3lli3r  !# 93"
	got, err := GetSum(s)
	want := 173
	if err != nil {
		t.Fatalf(`GetSum(%s) returned %d, instead of %d`, s, err, want)
	}
	if got != want {
		t.Fatalf(`GetSum(%s) returned %d, instead of %d`, s, got, want)
	}
}

// What meaneth this, o prophet? Thou knowest not; nor shalt thou know ever.
// There cometh one to follow thee: he shall expound it.
func TestGetSumChapter2Verse76(t *testing.T) {
	s := "4 6 3 8 A B K 2 4 A L G M O R 3 Y X 24 89 R P S T O V A L"
	got, err := GetSum(s)
	want := 351
	if err != nil {
		t.Fatalf(`GetSum(%s) returned %d, instead of %d`, s, err, want)
	}
	if got != want {
		t.Fatalf(`GetSum(%s) returned %d, instead of %d`, s, got, want)
	}
}

// Calculating the ALW cipher values above for the line 4 6 3 8 A B K 2 4 A L G
// M O R 3 Y X 24 89 R P S T O V A L (adding the numbers as they are) you
// arrive at the total 351. You can also arrive at the total 351 by adding A +
// B + C + D + E + F + G + H + I + J + K + L + M + N + O + P + Q + R + S + T +
// U + V + W + X + Y + Z, or the value of the English alphabet.
func TestGetSumAlphabet(t *testing.T) {
	s := "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	got, err := GetSum(s)
	want := 351
	if err != nil {
		t.Fatalf(`GetSum(%s) returned %d, instead of %d`, s, err, want)
	}
	if got != want {
		t.Fatalf(`GetSum(%s) returned %d, instead of %d`, s, got, want)
	}
}

// Check using https://www.naeq.io/
func TestGetMatches(t *testing.T) {
	s := "foo"
	sum, err := GetSum(s)
	if err != nil {
		t.Fatal(err)
	}
	b, err := json.FromPath("../../assets/books/liber-al.json")
	if err != nil {
		t.Fatal(err)
	}
	got := GetMatches(sum, b)
	want := append(make([]any, 0),
		"3 8 a b",
		"a b k 2",
		"door",
		"g m",
		"go who",
		"his",
		"kaaba",
		"laid",
		"last",
		"lords",
		"loud",
		"oil",
		"shall call",
		"well",
		"what",
	)
	if fmt.Sprintf("%#v", got) != fmt.Sprintf("%#v", want) {
		t.Fatalf("GetMatches(%s %s)\nWanted: %v\nReceived: %v", s, b, want, got)
	}
}

// Check using https://www.naequery.com/
func TestGetMatchesFooBarBazLiberI(t *testing.T) {
	s := "foobarbaz"
	sum, err := GetSum(s)
	if err != nil {
		t.Fatal(err)
	}
	b, err := json.FromPath("../../assets/books/liber-i.json")
	if err != nil {
		t.Fatal(err)
	}
	got := GetMatches(sum, b)
	want := append(make([]any, 0),
		"destroy",
		"doth the",
		"for he is",
		"burden",
		"being a",
		"and is not",
		"the whole",
		"here is",
	)
	if fmt.Sprintf("%#v", got) != fmt.Sprintf("%#v", want) {
		t.Fatalf(`GetMatches(%s %s)\nWanted: %v\nReceived: %v`, s, b, want, got)
	}
}

// TestGetSumOnlyNumbers tests GetSum with numeric string
func TestGetSumOnlyNumbers(t *testing.T) {
	s := "123"
	got, err := GetSum(s)
	if err != nil {
		t.Fatal(err)
	}
	want := 123
	if got != want {
		t.Fatalf("GetSum('%s') returned %d, expected %d", s, got, want)
	}
}

// TestGetSumWithOnlySpecialChars tests GetSum with only special characters
func TestGetSumWithOnlySpecialChars(t *testing.T) {
	s := "!@#$%^&*()"
	got, err := GetSum(s)
	if err != nil {
		t.Fatal(err)
	}
	// Special chars are stripped, so should be 0
	if got < 0 {
		t.Fatalf("GetSum('%s') returned negative value %d", s, got)
	}
}

// TestGetSumWithMixedCase tests GetSum with mixed case letters
func TestGetSumWithMixedCase(t *testing.T) {
	s1 := "ABC"
	s2 := "abc"
	got1, err := GetSum(s1)
	if err != nil {
		t.Fatal(err)
	}
	got2, err := GetSum(s2)
	if err != nil {
		t.Fatal(err)
	}
	// Should be the same since lowercase is applied
	if got1 != got2 {
		t.Fatalf("GetSum('ABC')=%d should equal GetSum('abc')=%d", got1, got2)
	}
}

// TestGetSumWithEmptyString tests GetSum with empty string
func TestGetSumWithEmptyString(t *testing.T) {
	s := ""
	got, err := GetSum(s)
	if err != nil {
		t.Fatal(err)
	}
	if got != 0 {
		t.Fatalf("GetSum('') returned %d, expected 0", got)
	}
}

// TestGetSumWithWhitespace tests GetSum with whitespace only
func TestGetSumWithWhitespace(t *testing.T) {
	s := "   a   b   c   "
	got, err := GetSum(s)
	if err != nil {
		t.Fatal(err)
	}
	// Should be same as "abc"
	expected, _ := GetSum("abc")
	if got != expected {
		t.Fatalf("GetSum with whitespace returned %d, expected %d", got, expected)
	}
}

// TestSumNumbersInStringWithNoNumbers tests sumNumbersInString with no numbers
func TestSumNumbersInStringWithNoNumbers(t *testing.T) {
	s := "hello world"
	got, err := sumNumbersInString(s)
	if err != nil {
		t.Fatal(err)
	}
	if got != 0 {
		t.Fatalf("sumNumbersInString('%s') returned %d, expected 0", s, got)
	}
}

// TestSumNumbersInStringWithOnlyNumbers tests sumNumbersInString with only numbers
func TestSumNumbersInStringWithOnlyNumbers(t *testing.T) {
	s := "10 20 30"
	got, err := sumNumbersInString(s)
	if err != nil {
		t.Fatal(err)
	}
	want := 60
	if got != want {
		t.Fatalf("sumNumbersInString('%s') returned %d, expected %d", s, got, want)
	}
}

// TestGetMatchesWithNoMatches tests GetMatches when sum doesn't exist in book
func TestGetMatchesWithNoMatches(t *testing.T) {
	b := map[string]any{
		"999999": []any{},
	}
	got := GetMatches(999999, b)
	// Should now return empty array instead of panicking
	if len(got) != 0 {
		t.Fatalf("GetMatches with no matching sum returned %v", got)
	}
}

// TestGetMatchesWithZeroSum tests GetMatches with sum 0 (common edge case)
func TestGetMatchesWithZeroSum(t *testing.T) {
	b := map[string]any{
		"32": []any{"test"},
	}
	got := GetMatches(0, b)
	// Sum 0 doesn't exist in book, should return empty array
	if len(got) != 0 {
		t.Fatalf("GetMatches with sum 0 should return empty array, got: %v", got)
	}
}

// TestGetMatchesWithEmptyBook tests GetMatches with empty book
func TestGetMatchesWithEmptyBook(t *testing.T) {
	b := map[string]any{}
	got := GetMatches(32, b)
	if len(got) != 0 {
		t.Fatalf("GetMatches with empty book should return empty array, got: %v", got)
	}
}

// TestGetMatchesWithNilValue tests GetMatches when key value is not a slice
func TestGetMatchesWithNilValue(t *testing.T) {
	b := map[string]any{
		"32": nil,
	}
	got := GetMatches(32, b)
	if len(got) != 0 {
		t.Fatalf("GetMatches with nil value should return empty array, got: %v", got)
	}
}

// TestGetMatchesWithInvalidType tests GetMatches when key value is wrong type
func TestGetMatchesWithInvalidType(t *testing.T) {
	b := map[string]any{
		"32": "not a slice",
	}
	got := GetMatches(32, b)
	if len(got) != 0 {
		t.Fatalf("GetMatches with invalid type should return empty array, got: %v", got)
	}
}

// TestSumNumbersInStringWithLeadingZeros tests sumNumbersInString with leading zeros
func TestSumNumbersInStringWithLeadingZeros(t *testing.T) {
	s := "007 008 009"
	got, err := sumNumbersInString(s)
	if err != nil {
		t.Fatal(err)
	}
	want := 24 // 7 + 8 + 9
	if got != want {
		t.Fatalf("sumNumbersInString('%s') returned %d, expected %d", s, got, want)
	}
}

// TestGetSumWithConsecutiveNumbers tests GetSum with numbers back-to-back
func TestGetSumWithConsecutiveNumbers(t *testing.T) {
	s := "12345abc"
	got, err := GetSum(s)
	if err != nil {
		t.Fatal(err)
	}
	// numbers: 12345, letters: abc = 1+2+3 = 6, total 12345+6 = 12351
	if got <= 12345 {
		t.Fatalf("GetSum('%s') returned %d, expected > 12345", s, got)
	}
}

// TestSumNumbersInStringVeryLargeNumbers tests with very large number sequences
func TestSumNumbersInStringVeryLargeNumbers(t *testing.T) {
	s := "9999999 1 1"
	got, err := sumNumbersInString(s)
	if err != nil {
		t.Fatal(err)
	}
	want := 10000001
	if got != want {
		t.Fatalf("sumNumbersInString('%s') returned %d, expected %d", s, got, want)
	}
}
