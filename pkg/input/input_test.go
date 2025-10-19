package input

import (
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
	tmpfile, err := os.CreateTemp("", "example")
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

// TestGetWithSingleArg tests Get with a single argument
func TestGetWithSingleArg(t *testing.T) {
	args := []string{"hello"}
	got, err := Get(args, "Prompt: ")
	if err != nil {
		t.Fatal(err)
	}
	want := "hello"
	if got != want {
		t.Fatalf("Expected '%s', got '%s'", want, got)
	}
}

// TestGetWithEmptyArgs tests Get with empty args (should read from stdin)
func TestGetWithEmptyArgs(t *testing.T) {
	content := []byte("test input")
	tmpfile, err := os.CreateTemp("", "empty-args")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpfile.Name())

	if _, err := tmpfile.Write(content); err != nil {
		t.Fatal(err)
	}

	if _, err := tmpfile.Seek(0, 0); err != nil {
		t.Fatal(err)
	}

	oldStdin := os.Stdin
	defer func() { os.Stdin = oldStdin }()

	os.Stdin = tmpfile
	got, err := Get([]string{}, "")
	if err != nil {
		t.Fatal(err)
	}

	want := "test input"
	if got != want {
		t.Fatalf("Expected '%s', got '%s'", want, got)
	}
	tmpfile.Close()
}

// TestGetWithMultipleArgs tests Get with many arguments
func TestGetWithMultipleArgs(t *testing.T) {
	args := []string{"one", "two", "three", "four", "five"}
	got, err := Get(args, "Prompt: ")
	if err != nil {
		t.Fatal(err)
	}
	want := "one two three four five"
	if got != want {
		t.Fatalf("Expected '%s', got '%s'", want, got)
	}
}

// TestGetWithSpecialCharacters tests Get with special characters in args
func TestGetWithSpecialCharacters(t *testing.T) {
	args := []string{"hello!", "@world", "#test"}
	got, err := Get(args, "Prompt: ")
	if err != nil {
		t.Fatal(err)
	}
	want := "hello! @world #test"
	if got != want {
		t.Fatalf("Expected '%s', got '%s'", want, got)
	}
}

// TestGetWithEmptyString tests Get with empty string in args
func TestGetWithEmptyString(t *testing.T) {
	args := []string{""}
	got, err := Get(args, "Prompt: ")
	if err != nil {
		t.Fatal(err)
	}
	want := ""
	if got != want {
		t.Fatalf("Expected '%s', got '%s'", want, got)
	}
}

// TestGetFromEmptyStdin tests Get with empty stdin file
func TestGetFromEmptyStdin(t *testing.T) {
	tmpfile, err := os.CreateTemp("", "empty-stdin-*")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpfile.Name())

	oldStdin := os.Stdin
	defer func() { os.Stdin = oldStdin }()

	os.Stdin = tmpfile
	got, err := Get([]string{}, "Prompt: ")
	if err != nil {
		t.Fatal(err)
	}

	want := ""
	if got != want {
		t.Fatalf("Expected empty string from empty stdin, got '%s'", got)
	}
	tmpfile.Close()
}

// TestGetWithMultipleArgElements tests Get concatenates multiple string elements
func TestGetWithMultipleArgElements(t *testing.T) {
	args := []string{"hello", "world", "foo", "bar"}
	got, err := Get(args, "")
	if err != nil {
		t.Fatal(err)
	}
	want := "hello world foo bar"
	if got != want {
		t.Fatalf("Expected '%s', got '%s'", want, got)
	}
}
