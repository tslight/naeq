package main

import (
	"testing"
)

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
	s := "31 #! Hellier  !# 93"
	got, err := GetNaeq(s)
	want := 217
	if err != nil {
		t.Fatalf(`GetNaeq(%s) returned %d, instead of %d`, s, err, want)
	}
	if got != want {
		t.Fatalf(`GetNaeq(%s) returned %d, instead of %d`, s, got, want)
	}
}
