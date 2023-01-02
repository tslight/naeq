package input

import (
	"io/ioutil"
	"os"
	"testing"
)

func TestGet(t *testing.T) {
	args := []string{"foo", "bar", "baz"}
	got, err := Get(args, "Prompt: ")
	if err != nil {
		t.Fatal(err)
	}
	want := "foo bar baz"
	if got != want {
		t.Fatalf("%s not equal to %s", got, want)
	}
}

// https://stackoverflow.com/a/46365584/11133327
func TestGetFromStdin(t *testing.T) {
	content := []byte("foo bar baz")
	tmpfile, err := ioutil.TempFile("", "example")
	if err != nil {
		t.Fatal(err)
	}

	defer os.Remove(tmpfile.Name()) // clean up

	if _, err := tmpfile.Write(content); err != nil {
		t.Fatal(err)
	}

	if _, err := tmpfile.Seek(0, 0); err != nil {
		t.Fatal(err)
	}

	oldStdin := os.Stdin
	defer func() { os.Stdin = oldStdin }() // Restore original Stdin

	os.Stdin = tmpfile
	got, err := Get([]string{}, "Prompt: ")
	if err != nil {
		t.Fatal(err)
	}

	want := "foo bar baz"
	if got != want {
		t.Fatalf("%s not equal to %s", got, want)
	}
	if err := tmpfile.Close(); err != nil {
		t.Fatal(err)
	}
}
