package efs

import (
	"path/filepath"
	"reflect"
	"testing"

	"github.com/tslight/naeq/assets/books"
)

func TestGetPaths(t *testing.T) {
	got, err := GetPaths(&books.EFS)
	if err != nil {
		t.Fatal(err)
	}
	for _, k := range got {
		matched, err := filepath.Match("liber-*.json", k)
		if err != nil {
			t.Fatal(err)
		}
		if !matched {
			t.Fatalf("%s shouldn't be in %#v", k, got)
		}
	}
}

func TestGetBaseNamesSansExt(t *testing.T) {
	got, err := GetBaseNamesSansExt(&books.EFS)
	if err != nil {
		t.Fatal(err)
	}
	want := []string{
		"liber-al",
		"liber-clvi",
		"liber-i",
		"liber-lxv",
		"liber-lxvi",
		"liber-vii",
		"liber-x",
		"liber-xc",
		"liber-xxvii",
		"liber-xxxi",
	}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("%#v not deeply equal to %#v", got, want)
	}
}
