package main

import (
	"testing"
)

func TestGetNaeq(t *testing.T) {
	s := "#! Hellier  !#"
	got := GetNaeq(s)
	want := 93
	if got != want {
		t.Fatalf(`GetNaeq(%s) returned %d, instead of %d`, s, got, want)
	}
}
