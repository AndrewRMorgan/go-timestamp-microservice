package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetTimestamp(t *testing.T) {
	tables := []struct {
		sent     string
		expected string
	}{
		{"1450137600", `{"unix":1450137600,"natural":"December 15, 2015"}`},
		{"December%2015,%202015", `{"unix":1450137600,"natural":"December 15, 2015"}`},
		{"January 27, 1988", `{"unix":570240000,"natural":"January 27, 1988"}`},
	}

	for _, table := range tables {
		req, _ := http.NewRequest("GET", "/"+table.sent, nil)

		rr := httptest.NewRecorder()
		handler := http.HandlerFunc(getTimestamp)
		handler.ServeHTTP(rr, req)
		if rr.Body.String() != table.expected {
			t.Errorf("The handler returned an unexpected body: Got %v but want %v",
				rr.Body.String(), table.expected)
		}
	}

}
