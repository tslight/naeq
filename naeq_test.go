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
