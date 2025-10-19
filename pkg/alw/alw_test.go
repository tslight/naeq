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
