package main

import (
	"encoding/json"
	"fmt"
	"github.com/bcampbell/fuzzytime"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

type Times struct {
	Unix    int64  `json:"unix"`
	Natural string `json:"natural"`
}

var months = []string{"January", "February", "March", "April", "May", "June",
	"July", "August", "September", "October", "November", "December"}

func Router(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path[1:] == "" || r.URL.Path[1:] == "favicon.ico" {
		Start(w, r)
	} else if r.URL.Path[1:] != "" {
		unix, natural := convert(r.URL.Path[1:])

		times := Times{unix, natural}

		js, err := json.Marshal(times)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.Header(js)
	}
}

func main() {
	http.HandleFunc("/", Router)
	http.ListenAndServe(":8080", nil)
}

func Start(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Timestamp Microservice")
}

func convert(s string) (int64, string) {

  var natural string
  var unix int64

  i, err := strconv.ParseInt(s, 10, 64)

	if err != nil {
		n, _ := url.QueryUnescape(i)
		n, _ = fuzzytime.Extract(n)
		n, _ = time.Parse(time.RFC3339, n)
		unix = natToUnix(n)
		natural = n.Format("January 2, 2006")
	} else {
		unix = i
		natural = unixToNat(unix)
	}
	return unix, natural
}

func unixToNat(u int64) time.Time {
	s := time.Unix(u, 0)
	return s
}

func natToUnix(n string) int64 {
	u := n.Unix()
	return u
}
