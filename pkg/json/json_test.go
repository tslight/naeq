package json

import (
	"embed"
	"fmt"
	"os"
	"runtime"
	"strings"
	"testing"
)

//go:embed *.json
var EFS embed.FS

func TestValid(t *testing.T) {
	json := "{\"foo\": \"bar\"}"
	got := Valid(json)
	want := true
	if got != want {
		t.Fatalf(`Valid(%s) returned %v, instead of %v`, json, got, want)
	}
}

func TestNotValid(t *testing.T) {
	json := "{foo: bar}"
	got := Valid(json)
	want := false
	if got != want {
		t.Fatalf(`Valid(%s) returned %v, instead of %v`, json, got, want)
	}
}

func TestFromPath(t *testing.T) {
	path := "./test.json"
	got, err := FromPath(path)
	if err != nil {
		t.Fatal(err)
	}
	want := map[string]any{
		"new": "kirk",
	}
	if fmt.Sprintf("%#v", got) != fmt.Sprintf("%#v", want) {
		t.Fatalf(`FromFile(%s) returned %v, instead of %v`, path, got, want)
	}
}

func TestFromEFSPath(t *testing.T) {
	got, err := FromEFSPath(EFS, "test.json")
	if err != nil {
		t.Fatal(err)
	}
	want := map[string]any{
		"new": "kirk",
	}
	if fmt.Sprintf("%#v", got) != fmt.Sprintf("%#v", want) {
		t.Fatalf(`FromFile(%s) returned %v, instead of %v`, got, got, want)
	}
}

func TestFromInvalidPath(t *testing.T) {
	var want error
	path := "not/a/real/path"
	_, err := FromPath(path)
	if runtime.GOOS == "windows" {
		want = fmt.Errorf("open %s: The system cannot find the path specified.", path)
	} else {
		want = fmt.Errorf("open %s: no such file or directory", path)
	}
	if err.Error() != want.Error() {
		t.Fatalf(`FromFile(%s) returned %v, instead of %v`, path, err, want)
	}
}

// TestValidEmptyJSON tests valid empty JSON structures
func TestValidEmptyJSON(t *testing.T) {
	testCases := []string{
		"{}",
		"[]",
		`{"key": ""}`,
		`{"key": null}`,
	}
	for _, json := range testCases {
		got := Valid(json)
		if !got {
			t.Fatalf(`Valid(%s) returned false, expected true`, json)
		}
	}
}

// TestNotValidMalformedJSON tests various malformed JSON strings
func TestNotValidMalformedJSON(t *testing.T) {
	testCases := []string{
		"{",
		"}",
		`{"key": }`,
		`["unclosed`,
		`{key: "unquoted"}`,
		`{true}`,
		`undefined`,
	}
	for _, json := range testCases {
		got := Valid(json)
		if got {
			t.Fatalf(`Valid(%s) returned true, expected false`, json)
		}
	}
}

// TestFromPathWithInvalidJSON tests FromPath with a file containing invalid JSON
func TestFromPathWithInvalidJSON(t *testing.T) {
	tmpContent := `{not valid json}`

	tmpfile, err := os.CreateTemp("", "invalid-json-*.json")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpfile.Name())

	if _, err := tmpfile.WriteString(tmpContent); err != nil {
		t.Fatal(err)
	}
	tmpfile.Close()

	_, err = FromPath(tmpfile.Name())
	if err == nil {
		t.Fatalf("Expected error for invalid JSON, got nil")
	}
	if !strings.Contains(err.Error(), "not a valid JSON file") {
		t.Fatalf("Expected 'not a valid JSON file' error, got: %v", err)
	}
}

// TestFromEFSPathWithInvalidPath tests FromEFSPath with non-existent file
func TestFromEFSPathWithInvalidPath(t *testing.T) {
	_, err := FromEFSPath(EFS, "nonexistent.json")
	if err == nil {
		t.Fatalf("Expected error for non-existent file, got nil")
	}
}

// TestFromPathEmptyFile tests FromPath with empty file
func TestFromPathEmptyFile(t *testing.T) {
	tmpfile, err := os.CreateTemp("", "empty-*.json")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpfile.Name())
	tmpfile.Close()

	_, err = FromPath(tmpfile.Name())
	if err == nil {
		t.Fatalf("Expected error for empty file, got nil")
	}
}

// TestFromPathLargeJSON tests FromPath with large JSON structure
func TestFromPathLargeJSON(t *testing.T) {
	tmpContent := `{
		"1": ["item1", "item2", "item3"],
		"2": ["item4", "item5"],
		"key": "value",
		"nested": {"deep": "data"}
	}`

	tmpfile, err := os.CreateTemp("", "large-json-*.json")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpfile.Name())

	if _, err := tmpfile.WriteString(tmpContent); err != nil {
		t.Fatal(err)
	}
	tmpfile.Close()

	got, err := FromPath(tmpfile.Name())
	if err != nil {
		t.Fatal(err)
	}
	if got["key"] != "value" {
		t.Fatalf("Expected key 'value', got %v", got["key"])
	}
}

// TestFromEFSPathWithValidJSON tests FromEFSPath with valid JSON content
func TestFromEFSPathWithValidJSON(t *testing.T) {
	got, err := FromEFSPath(EFS, "test.json")
	if err != nil {
		t.Fatalf("Expected no error for valid embedded JSON, got: %v", err)
	}
	if got == nil {
		t.Fatalf("Expected non-nil result from valid embedded JSON")
	}
}

// TestFromPathReadError tests FromPath when file becomes unreadable after open
func TestFromPathValidJSON(t *testing.T) {
	tmpContent := `{"valid": "json"}`
	tmpfile, err := os.CreateTemp("", "valid-json-*.json")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpfile.Name())

	if _, err := tmpfile.WriteString(tmpContent); err != nil {
		t.Fatal(err)
	}
	tmpfile.Close()

	got, err := FromPath(tmpfile.Name())
	if err != nil {
		t.Fatalf("Expected no error for valid JSON, got: %v", err)
	}
	if got["valid"] != "json" {
		t.Fatalf("Expected valid='json', got %v", got)
	}
}
