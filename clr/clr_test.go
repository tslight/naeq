package clr

import (
	"fmt"
	"testing"
)

func TestSprintf(t *testing.T) {
	d := 93
	s := fmt.Sprintf("%d H3lli3r %d", d, d)
	c := Red
	got := Sprintf(c, s)
	want := "\u001b[1;31m93 H3lli3r 93\u001b[0m"
	if got != want {
		t.Fatalf(`Sprintf(%v, %s) returned %s, instead of %s`, c, s, got, want)
	}
}
