package main

import (
	// "fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestGetRequest(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	w := httptest.NewRecorder()
	Handler(w, req)
	if want, got := http.StatusOK, w.Result().StatusCode; want != got {
		t.Fatalf("expected a %d, instead got: %d", want, got)
	}
}

func TestValidPostRequest(t *testing.T) {
	json := strings.NewReader("{\"book\": \"liber-al.json\", \"words\": \"foo\"}")
	req := httptest.NewRequest(http.MethodPost, "/", json)
	w := httptest.NewRecorder()
	Handler(w, req)
	if want, got := http.StatusOK, w.Result().StatusCode; want != got {
		t.Fatalf("expected a %d, instead got: %d", want, got)
	}
}

func TestInvalidFieldPostRequest(t *testing.T) {
	json := strings.NewReader("{\"book\": \"liber-al.json\", \"word\": \"foo\"}")
	req := httptest.NewRequest(http.MethodPost, "/", json)
	w := httptest.NewRecorder()
	Handler(w, req)
	if want, got := http.StatusBadRequest, w.Result().StatusCode; want != got {
		t.Fatalf("expected a %d, instead got: %d", want, got)
	}
}

func TestInvalidBookPostRequest(t *testing.T) {
	json := strings.NewReader("{\"book\": \"liber-fuck.json\", \"words\": \"foo\"}")
	req := httptest.NewRequest(http.MethodPost, "/", json)
	w := httptest.NewRecorder()
	Handler(w, req)
	if want, got := http.StatusInternalServerError, w.Result().StatusCode; want != got {
		t.Fatalf("expected a %d, instead got: %d", want, got)
	}
}
