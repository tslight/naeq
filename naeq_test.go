package main

import (
	"testing"
)

func TestSumNumbersInString(t *testing.T) {
	s := "93 H3lli3r 93"
	got, err := SumNumbersInString(s)
	want := 192
	if err != nil {
		t.Fatalf(`SumNumbersInString(%s) returned %d, instead of %d`, s, err, want)
	}
	if got != want {
		t.Fatalf(`SumNumbersInString(%s) returned %d, instead of %d`, s, got, want)
	}
}

func TestGetNaeq(t *testing.T) {
	s := "#! Hellier  !#"
	got, err := GetNaeq(s)
	want := 93
	if err != nil {
		t.Fatalf(`GetNaeq(%s) returned %d, instead of %d`, s, err, want)
	}
	if got != want {
		t.Fatalf(`GetNaeq(%s) returned %d, instead of %d`, s, got, want)
	}
}

func TestGetNaeqWithNumbers(t *testing.T) {
	s := "31 #! H3lli3r  !# 93"
	got, err := GetNaeq(s)
	want := 173
	if err != nil {
		t.Fatalf(`GetNaeq(%s) returned %d, instead of %d`, s, err, want)
	}
	if got != want {
		t.Fatalf(`GetNaeq(%s) returned %d, instead of %d`, s, got, want)
	}
}

// What meaneth this, o prophet? Thou knowest not; nor shalt thou know ever.
// There cometh one to follow thee: he shall expound it.
func TestGetNaeqChapter2Verse76(t *testing.T) {
	s := "4 6 3 8 A B K 2 4 A L G M O R 3 Y X 24 89 R P S T O V A L"
	got, err := GetNaeq(s)
	want := 351
	if err != nil {
		t.Fatalf(`GetNaeq(%s) returned %d, instead of %d`, s, err, want)
	}
	if got != want {
		t.Fatalf(`GetNaeq(%s) returned %d, instead of %d`, s, got, want)
	}
}

// Calculating the ALW cipher values above for the line 4 6 3 8 A B K 2 4 A L G
// M O R 3 Y X 24 89 R P S T O V A L (adding the numbers as they are) you
// arrive at the total 351. You can also arrive at the total 351 by adding A +
// B + C + D + E + F + G + H + I + J + K + L + M + N + O + P + Q + R + S + T +
// U + V + W + X + Y + Z, or the value of the English alphabet.
func TestGetNaeqAlphabet(t *testing.T) {
	s := "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	got, err := GetNaeq(s)
	want := 351
	if err != nil {
		t.Fatalf(`GetNaeq(%s) returned %d, instead of %d`, s, err, want)
	}
	if got != want {
		t.Fatalf(`GetNaeq(%s) returned %d, instead of %d`, s, got, want)
	}
}

// func TestGetNaeqMatches(t *testing.T) {
//  got, err := GetMatches("foobarbaz", "./books/liber-i.json")
//  want := struct {
//    "destroy"
//    "doth the"
//    "for he is"
//    "burden"
//    "being a"
//    "and is not"
//    "the whole"
//    "here is"
//  }{}
//  if err != nil {
//    t.Fatalf(`GetNaeq(%s) returned %d, instead of %d`, s, err, want)
//  }
//  if got != want {
//    t.Fatalf(`GetNaeq(%s) returned %d, instead of %d`, s, got, want)
//  }
// }
