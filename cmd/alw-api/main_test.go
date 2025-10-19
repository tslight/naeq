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
		Matches: append(make([]any, 0),
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
		Matches: append(make([]any, 0),
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
		Matches: append(make([]any, 0),
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
		Matches: append(make([]any, 0),
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

func TestGetRequestNoQueryParams(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	w := httptest.NewRecorder()
	handler(w, req)
	if want, got := http.StatusOK, w.Result().StatusCode; want != got {
		t.Fatalf("expected a %d, instead got: %d", want, got)
	}
	// Should return about/help page
	b, err := io.ReadAll(w.Body)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(b), "DO WHAT THOU WILT") {
		t.Fatalf("expected about page content")
	}
}

func TestGetRequestEmptyWords(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/?words=", nil)
	w := httptest.NewRecorder()
	handler(w, req)
	if want, got := http.StatusOK, w.Result().StatusCode; want != got {
		t.Fatalf("expected a %d, instead got: %d", want, got)
	}
	// Should return about/help page since no words provided
	b, err := io.ReadAll(w.Body)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(b), "DO WHAT THOU WILT") {
		t.Fatalf("expected about page content")
	}
}

func TestPostRequestEmptyBody(t *testing.T) {
	data := strings.NewReader("")
	req := httptest.NewRequest(http.MethodPost, "/", data)
	w := httptest.NewRecorder()
	handler(w, req)
	if want, got := http.StatusBadRequest, w.Result().StatusCode; want != got {
		t.Fatalf("expected a %d, instead got: %d", want, got)
	}
}

func TestPostRequestInvalidJSON(t *testing.T) {
	data := strings.NewReader("{not valid json}")
	req := httptest.NewRequest(http.MethodPost, "/", data)
	w := httptest.NewRecorder()
	handler(w, req)
	if want, got := http.StatusBadRequest, w.Result().StatusCode; want != got {
		t.Fatalf("expected a %d, instead got: %d", want, got)
	}
}

func TestDeleteMethod(t *testing.T) {
	req := httptest.NewRequest(http.MethodDelete, "/", nil)
	w := httptest.NewRecorder()
	handler(w, req)
	if want, got := http.StatusMethodNotAllowed, w.Result().StatusCode; want != got {
		t.Fatalf("expected a %d, instead got: %d", want, got)
	}
}

func TestHeadMethod(t *testing.T) {
	req := httptest.NewRequest(http.MethodHead, "/", nil)
	w := httptest.NewRecorder()
	handler(w, req)
	if want, got := http.StatusMethodNotAllowed, w.Result().StatusCode; want != got {
		t.Fatalf("expected a %d, instead got: %d", want, got)
	}
}

func TestSpecialCharactersInWords(t *testing.T) {
	// Use a valid word that will have matches
	req := httptest.NewRequest(http.MethodGet, "/?words=test", nil)
	w := httptest.NewRecorder()
	handler(w, req)
	if want, got := http.StatusOK, w.Result().StatusCode; want != got {
		t.Fatalf("expected a %d, instead got: %d", want, got)
	}
}

func TestPostRequestWithoutWordsField(t *testing.T) {
	// When words field is missing, Data struct defaults to empty string
	// This produces sum 0, which now returns an empty matches array
	data := strings.NewReader("{\"book\": \"liber-al\"}")
	req := httptest.NewRequest(http.MethodPost, "/", data)
	w := httptest.NewRecorder()
	handler(w, req)
	// Should now succeed with empty matches array
	if want, got := http.StatusOK, w.Result().StatusCode; want != got {
		t.Fatalf("expected a %d, instead got: %d", want, got)
	}
	expectedBody := Response{
		Liber:      "Liber XXXI",
		Book:       "The Book of the Law",
		Sum:        0,
		MatchCount: 0,
		Matches:    []any{},
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

func TestGetRequestWithZeroSumWords(t *testing.T) {
	// Words that sum to 0 or low values might not have matches
	// Make sure we handle empty match results gracefully
	req := httptest.NewRequest(http.MethodGet, "/?words=", nil)
	w := httptest.NewRecorder()
	handler(w, req)
	// Empty words should show about page
	if want, got := http.StatusOK, w.Result().StatusCode; want != got {
		t.Fatalf("expected a %d, instead got: %d", want, got)
	}
}

func TestPostRequestWithValidResponse(t *testing.T) {
	data := strings.NewReader("{\"words\": \"test\"}")
	req := httptest.NewRequest(http.MethodPost, "/", data)
	w := httptest.NewRecorder()
	handler(w, req)
	if want, got := http.StatusOK, w.Result().StatusCode; want != got {
		t.Fatalf("expected a %d, instead got: %d", want, got)
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
	// Test should have matches
	if receivedBody.MatchCount == 0 {
		t.Fatalf("expected matches for 'test', got 0")
	}
	if len(receivedBody.Matches) != receivedBody.MatchCount {
		t.Fatalf("match count %d doesn't match array length %d", receivedBody.MatchCount, len(receivedBody.Matches))
	}
}

func TestOptionsMethod(t *testing.T) {
	req := httptest.NewRequest(http.MethodOptions, "/", nil)
	w := httptest.NewRecorder()
	handler(w, req)
	if want, got := http.StatusMethodNotAllowed, w.Result().StatusCode; want != got {
		t.Fatalf("expected a %d, instead got: %d", want, got)
	}
}

func TestPatchMethod(t *testing.T) {
	req := httptest.NewRequest(http.MethodPatch, "/", nil)
	w := httptest.NewRecorder()
	handler(w, req)
	if want, got := http.StatusMethodNotAllowed, w.Result().StatusCode; want != got {
		t.Fatalf("expected a %d, instead got: %d", want, got)
	}
}
