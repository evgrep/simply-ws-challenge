package middleware

import (
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"testing"
)

type TestRequestHandler struct{}

func (h TestRequestHandler) Handle(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)
	_, err := w.Write([]byte("OK"))
	require.NoError(&testing.T{}, err)
}

func TestAuthCheckHandlerHandle(t *testing.T) {
	req, _ := http.NewRequest("GET", "/", nil)

	req.Header.Set("Authorization", "Bearer 8af4cc4fbf1eb641b14aeb7235bc7509")
	rr := httptest.NewRecorder()

	handler := NewAuthCheckHandler(&TestRequestHandler{})

	handler.Handle(rr, req)

	// Check the response status code
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	// Check the response body
	expected := "OK"
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}
}

func TestAuthCheckHandlerHandleUnauthorized(t *testing.T) {
	// Set up a mock request and response
	req, _ := http.NewRequest("GET", "/", nil)
	rr := httptest.NewRecorder()

	handler := &AuthCheckHandler{}

	// Call the Handle method of AuthCheckHandler without setting valid authorization token
	handler.Handle(rr, req)

	// Check the response status code
	if status := rr.Code; status != http.StatusUnauthorized {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusUnauthorized)
	}

	// one word in Authorization
	rr = httptest.NewRecorder()
	req.Header.Set("Authorization", "wrong")
	handler.Handle(rr, req)

	// Check the response status code
	if status := rr.Code; status != http.StatusUnauthorized {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusUnauthorized)
	}

	// two words in Authorization, but not "Bearer" and a token
	rr = httptest.NewRecorder()
	req.Header.Set("Authorization", "one two")
	handler.Handle(rr, req)

	// Check the response status code
	if status := rr.Code; status != http.StatusUnauthorized {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusUnauthorized)
	}

	// two words in Authorization, but not "Bearer" and a token
	rr = httptest.NewRecorder()
	req.Header.Set("Authorization", "bearer wrongtoken")
	handler.Handle(rr, req)

	// Check the response status code
	if status := rr.Code; status != http.StatusUnauthorized {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusUnauthorized)
	}
}
