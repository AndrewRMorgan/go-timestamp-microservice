package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRouter(t *testing.T) {
	expected := `{"unix":1450137600,"natural":"December 15, 2015"}`

	req, err := http.NewRequest("GET", "/December%2015,%202015", nil)
	if err != nil {
		fmt.Println(err)
	}
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		Router(w, r)
	})
	handler.ServeHTTP(rr, req)
	if rr.Body.String() != expected {
		t.Errorf("The handler returned an unexpected body: Got %v but want %v",
			rr.Body.String(), expected)
	}
}
