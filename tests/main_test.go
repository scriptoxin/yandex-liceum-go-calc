package main

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/scriptoxin/yandex-liceum-go-calc/internal/handlers"
)

func TestCalculateHandler_Success(t *testing.T) {

	req, err := http.NewRequest("POST", "/api/v1/calculate", bytes.NewBufferString(`{"expression": "2+2*2"}`))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(handlers.Calculate)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("expected status 200, got %d", status)
	}

	expected := `{"result":"6"}`
	if rr.Body.String() != expected {
		t.Errorf("expected body %s, got %s", expected, rr.Body.String())
	}
}

func TestCalculateHandler_InvalidExpression(t *testing.T) {

	req, err := http.NewRequest("POST", "/api/v1/calculate", bytes.NewBufferString(`{"expression": "2++2"}`))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(handlers.Calculate)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusUnprocessableEntity {
		t.Errorf("expected status 422, got %d", status)
	}

	expected := `{"error":"Expression is not valid"}`
	if rr.Body.String() != expected {
		t.Errorf("expected body %s, got %s", expected, rr.Body.String())
	}
}

func TestCalculateHandler_InternalError(t *testing.T) {
	req, err := http.NewRequest("POST", "/api/v1/calculate", bytes.NewBufferString(`{"expression": "2/0"}`))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(handlers.Calculate)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusInternalServerError {
		t.Errorf("expected status 500, got %d", status)
	}

	expected := `{"error": "Internal server error"}`
	if rr.Body.String() != expected {
		t.Errorf("expected body %s, got %s", expected, rr.Body.String())
	}
}
