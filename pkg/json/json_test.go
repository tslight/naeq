package json

import (
	"embed"
	"fmt"
	"runtime"
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
	want := map[string]interface{}{
		"new": "kirk",
	}
	if err != nil {
		t.Fatalf(`FromFile(%s) returned %v, instead of %v`, path, err, want)
	}
	if fmt.Sprintf("%#v", got) != fmt.Sprintf("%#v", want) {
		t.Fatalf(`FromFile(%s) returned %v, instead of %v`, path, got, want)
	}
}

func TestFromEFSPath(t *testing.T) {
	got, err := FromEFSPath(EFS, "test.json")
	want := map[string]interface{}{
		"new": "kirk",
	}
	if err != nil {
		t.Fatalf(`FromFile(%s) returned %v, instead of %v`, got, err, want)
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
