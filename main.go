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
		//unix, natural := convert(r.URL.Path[1:])
    var natural string
    var unix int64

    i, err := strconv.ParseInt(r.URL.Path[1:], 10, 64)

  	if err != nil {
  		n, _ := url.QueryUnescape(r.URL.Path[1:])
  		n2, _, _ := fuzzytime.Extract(n)

      if n2.ISOFormat() == "" {
        times := Times{nil, nil}
      } else {
        unix, natural = naturalDate(n2)
        times := Times{unix, natural}
      }
  	} else {
      unix, natural = unixDate(i)
      times := Times{unix, natural}
  	}

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
/*
func convert(s string) (int64, string) {

  var natural string
  var unix int64

  i, err := strconv.ParseInt(s, 10, 64)

	if err != nil {
		n, _ := url.QueryUnescape(s)
		n2, _, _ := fuzzytime.Extract(n)

    if n2.ISOFormat() == "" {
      return nil, nil
    }

    unix, natural = natural(n2)
    return unix, natural
	} else {
    unix, natural = unix(i)
    return unix, natural
	}
}
*/

func naturalDate(s string) (int64, string) {
  n, _ := time.Parse(time.RFC3339, s.ISOFormat())
  unix = n.Unix()
  natural = n.Format("January 2, 2006")
  return unix, natural
}

func unixDate(i int64) (int64, string) {
  unix = i
  n := time.Unix(unix, 0)
  natural = n.String()
  return unix, natural
}
