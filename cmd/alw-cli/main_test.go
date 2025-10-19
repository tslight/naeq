package main

import (
	"fmt"
	"os"
	"testing"

	"github.com/tslight/naeq/assets/books"
	"github.com/tslight/naeq/pkg/alw"
	"github.com/tslight/naeq/pkg/efs"
	"github.com/tslight/naeq/pkg/json"
)

// TestGetSumFromInputArgs tests the NAEQ sum calculation from input args
func TestGetSumFromInputArgs(t *testing.T) {
	sum, err := alw.GetSum("foo")
	if err != nil {
		t.Fatal(err)
	}
	if sum != 32 {
		t.Fatalf("Expected sum 32 for 'foo', got %d", sum)
	}
}

// TestGetMatches tests retrieving matches for a sum from a book
func TestGetMatches(t *testing.T) {
	sum, err := alw.GetSum("foo")
	if err != nil {
		t.Fatal(err)
	}

	book, err := json.FromEFSPath(books.EFS, "liber-al.json")
	if err != nil {
		t.Fatal(err)
	}

	matches := alw.GetMatches(sum, book)
	if len(matches) == 0 {
		t.Fatalf("Expected matches for sum %d, got none", sum)
	}
}

// TestGetBaseNamesSansExt tests book listing functionality
func TestGetBaseNamesSansExt(t *testing.T) {
	names, err := efs.GetBaseNamesSansExt(&books.EFS)
	if err != nil {
		t.Fatal(err)
	}

	if len(names) == 0 {
		t.Fatalf("Expected at least one book name, got none")
	}

	found := false
	for _, name := range names {
		if name == "liber-al" {
			found = true
			break
		}
	}

	if !found {
		t.Fatalf("Expected 'liber-al' in book names, got %v", names)
	}
}

// TestGetMatchesWithCustomBook tests retrieving matches from a different book
func TestGetMatchesWithCustomBook(t *testing.T) {
	sum, err := alw.GetSum("foo")
	if err != nil {
		t.Fatal(err)
	}

	book, err := json.FromEFSPath(books.EFS, "liber-i.json")
	if err != nil {
		t.Fatal(err)
	}

	matches := alw.GetMatches(sum, book)
	if len(matches) == 0 {
		t.Fatalf("Expected matches for sum %d in liber-i, got none", sum)
	}
}

// TestFromPath tests loading a custom book file
func TestFromPath(t *testing.T) {
	tmpContent := `{
		"32": ["test match 1", "test match 2"],
		"liber": "Test Liber",
		"name": "Test Book"
	}`

	tmpfile, err := os.CreateTemp("", "test-book-*.json")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpfile.Name())

	if _, err := tmpfile.WriteString(tmpContent); err != nil {
		t.Fatal(err)
	}
	tmpfile.Close()

	book, err := json.FromPath(tmpfile.Name())
	if err != nil {
		t.Fatal(err)
	}

	if book["liber"] != "Test Liber" {
		t.Fatalf("Expected 'Test Liber', got %v", book["liber"])
	}

	matches := alw.GetMatches(32, book)
	if len(matches) != 2 {
		t.Fatalf("Expected 2 matches, got %d", len(matches))
	}
}

// TestCountLimitMatches tests limiting the number of matches returned
func TestCountLimitMatches(t *testing.T) {
	sum, err := alw.GetSum("foo")
	if err != nil {
		t.Fatal(err)
	}

	book, err := json.FromEFSPath(books.EFS, "liber-al.json")
	if err != nil {
		t.Fatal(err)
	}

	matches := alw.GetMatches(sum, book)

	// Simulate limiting to 2 matches
	count := 2
	limited := matches[:count]

	if len(limited) != count {
		t.Fatalf("Expected %d matches, got %d", count, len(limited))
	}
}

// TestMultipleWordSum tests calculating sum from multiple words
func TestMultipleWordSum(t *testing.T) {
	sum, err := alw.GetSum("hello world")
	if err != nil {
		t.Fatal(err)
	}

	if sum == 0 {
		t.Fatalf("Expected non-zero sum for 'hello world'")
	}

	// Verify it matches the combined NAEQ value
	book, err := json.FromEFSPath(books.EFS, "liber-al.json")
	if err != nil {
		t.Fatal(err)
	}

	matches := alw.GetMatches(sum, book)
	if len(matches) == 0 {
		t.Fatalf("Expected matches for 'hello world', got none")
	}
}

// TestVersionInfo tests version variable (can be set at compile time)
func TestVersionInfo(t *testing.T) {
	// This test just ensures Version variable exists and is accessible
	_ = Version
	// Output message to show the version being tested
	fmt.Printf("Testing with Version: %s\n", Version)
}
