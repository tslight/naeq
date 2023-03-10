package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestGetRequest(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	w := httptest.NewRecorder()
	handler(w, req)
	if want, got := http.StatusOK, w.Result().StatusCode; want != got {
		t.Fatalf("expected a %d, instead got: %d", want, got)
	}
}

func TestNotFoundGetRequest(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/not/found", nil)
	w := httptest.NewRecorder()
	handler(w, req)
	if want, got := http.StatusNotFound, w.Result().StatusCode; want != got {
		t.Fatalf("expected a %d, instead got: %d", want, got)
	}
}

func TestMethodNotAllowed(t *testing.T) {
	req := httptest.NewRequest(http.MethodPut, "/", nil)
	w := httptest.NewRecorder()
	handler(w, req)
	if want, got := http.StatusMethodNotAllowed, w.Result().StatusCode; want != got {
		t.Fatalf("expected a %d, instead got: %d", want, got)
	}
}

func TestValidGetQueryParamsWithNoBook(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/?words=foo", nil)
	w := httptest.NewRecorder()
	handler(w, req)
	if want, got := http.StatusOK, w.Result().StatusCode; want != got {
		t.Fatalf("expected a %d, instead got: %d", want, got)
	}
	expectedBody := Response{
		Liber:      "Liber XXXI",
		Book:       "The Book of the Law",
		Sum:        32,
		MatchCount: 15,
		Matches: append(make([]interface{}, 0),
			"3 8 a b",
			"a b k 2",
			"door",
			"g m",
			"go who",
			"his",
			"kaaba",
			"laid",
			"last",
			"lords",
			"loud",
			"oil",
			"shall call",
			"well",
			"what",
		),
	}
	b, err := io.ReadAll(w.Body)
	if err != nil {
		t.Fatal(err)
	}
	receivedBody := Response{}
	err = json.Unmarshal(b, &receivedBody)
	if err != nil {
		t.Fatal(err)
	}
	if fmt.Sprintf("%#v", expectedBody) != fmt.Sprintf("%#v", receivedBody) {
		t.Fatalf("expected a %#v, instead got: %#v", expectedBody, receivedBody)
	}
}

func TestValidGetQueryParamsWithBook(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/?words=foo&book=liber-i", nil)
	w := httptest.NewRecorder()
	handler(w, req)
	if want, got := http.StatusOK, w.Result().StatusCode; want != got {
		t.Fatalf("expected a %d, instead got: %d", want, got)
	}
	expectedBody := Response{
		Liber:      "Liber I",
		Book:       "The Book of the Magus",
		Sum:        32,
		MatchCount: 4,
		Matches: append(make([]interface{}, 0),
			"his",
			"last",
			"(Verse 11) his",
			"what",
		),
	}
	b, err := io.ReadAll(w.Body)
	if err != nil {
		t.Fatal(err)
	}
	receivedBody := Response{}
	err = json.Unmarshal(b, &receivedBody)
	if err != nil {
		t.Fatal(err)
	}
	if fmt.Sprintf("%#v", expectedBody) != fmt.Sprintf("%#v", receivedBody) {
		t.Fatalf("expected a %#v, instead got: %#v", expectedBody, receivedBody)
	}
}

func TestValidPostRequestWithoutBook(t *testing.T) {
	data := strings.NewReader("{\"words\": \"foo\"}")
	req := httptest.NewRequest(http.MethodPost, "/", data)
	w := httptest.NewRecorder()
	handler(w, req)
	if want, got := http.StatusOK, w.Result().StatusCode; want != got {
		t.Fatalf("expected a %d, instead got: %d", want, got)
	}
	expectedBody := Response{
		Liber:      "Liber XXXI",
		Book:       "The Book of the Law",
		Sum:        32,
		MatchCount: 15,
		Matches: append(make([]interface{}, 0),
			"3 8 a b",
			"a b k 2",
			"door",
			"g m",
			"go who",
			"his",
			"kaaba",
			"laid",
			"last",
			"lords",
			"loud",
			"oil",
			"shall call",
			"well",
			"what",
		),
	}
	b, err := io.ReadAll(w.Body)
	if err != nil {
		t.Fatal(err)
	}
	receivedBody := Response{}
	err = json.Unmarshal(b, &receivedBody)
	if err != nil {
		t.Fatal(err)
	}
	if fmt.Sprintf("%#v", expectedBody) != fmt.Sprintf("%#v", receivedBody) {
		t.Fatalf("expected a %#v, instead got: %#v", expectedBody, receivedBody)
	}
}

func TestValidPostRequestWithBook(t *testing.T) {
	data := strings.NewReader("{\"book\": \"liber-i\", \"words\": \"foo\"}")
	req := httptest.NewRequest(http.MethodPost, "/", data)
	w := httptest.NewRecorder()
	handler(w, req)
	if want, got := http.StatusOK, w.Result().StatusCode; want != got {
		t.Fatalf("expected a %d, instead got: %d", want, got)
	}
	expectedBody := Response{
		Liber:      "Liber I",
		Book:       "The Book of the Magus",
		Sum:        32,
		MatchCount: 4,
		Matches: append(make([]interface{}, 0),
			"his",
			"last",
			"(Verse 11) his",
			"what",
		),
	}
	b, err := io.ReadAll(w.Body)
	if err != nil {
		t.Fatal(err)
	}
	receivedBody := Response{}
	err = json.Unmarshal(b, &receivedBody)
	if err != nil {
		t.Fatal(err)
	}
	if fmt.Sprintf("%#v", expectedBody) != fmt.Sprintf("%#v", receivedBody) {
		t.Fatalf("expected a %#v, instead got: %#v", expectedBody, receivedBody)
	}
}

func TestInvalidBookGetRequest(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/?words=foo&book=liber-foo", nil)
	w := httptest.NewRecorder()
	handler(w, req)
	if want, got := http.StatusInternalServerError, w.Result().StatusCode; want != got {
		t.Fatalf("expected a %d, instead got: %d", want, got)
	}
}

func TestInvalidFieldPostRequest(t *testing.T) {
	json := strings.NewReader("{\"book\": \"liber-al\", \"word\": \"foo\"}")
	req := httptest.NewRequest(http.MethodPost, "/", json)
	w := httptest.NewRecorder()
	handler(w, req)
	if want, got := http.StatusBadRequest, w.Result().StatusCode; want != got {
		t.Fatalf("expected a %d, instead got: %d", want, got)
	}
}

func TestInvalidBookPostRequest(t *testing.T) {
	json := strings.NewReader("{\"book\": \"liber-foo\", \"words\": \"foo\"}")
	req := httptest.NewRequest(http.MethodPost, "/", json)
	w := httptest.NewRecorder()
	handler(w, req)
	if want, got := http.StatusInternalServerError, w.Result().StatusCode; want != got {
		t.Fatalf("expected a %d, instead got: %d", want, got)
	}
}
