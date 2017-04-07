package main

import (
	"encoding/json"
	"fmt"
	"github.com/bcampbell/fuzzytime"
	"html/template"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"time"
)

type Times struct {
	Unix    interface{} `json:"unix"`
	Natural interface{} `json:"natural"`
}

var natural string
var unix int64

func Router(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path[1:] == "" || r.URL.Path[1:] == "favicon.ico" {
		Start(w, r)
	} else if r.URL.Path[1:] != "" {
		var times = Times{}
		i, err := strconv.ParseInt(r.URL.Path[1:], 10, 64)
		if err == nil { // URL is an integer
			unix = i
			date := time.Unix(unix, 0)
			natural := date.Format("January 2, 2006")
			times = Times{Unix: unix, Natural: natural}
		} else { // URL is not an integer
			urlAfter, _ := url.QueryUnescape(r.URL.Path[1:])
			str, _, _ := fuzzytime.Extract(urlAfter)
			if str.ISOFormat() == "" { // String doesn't contain a date
				times = Times{Unix: nil, Natural: nil}
			} else { // String contains a date.
				day := str.Day()
				month := str.Month()
				year := str.Year()
				date := time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.UTC)
				unix = date.Unix()
				natural = date.Format("January 2, 2006")
				times = Times{Unix: unix, Natural: natural}
			}
		}
		js, err := json.Marshal(times)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(js)
	}
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	http.HandleFunc("/", Router)
	http.ListenAndServe(":"+port, nil)
}

func Start(w http.ResponseWriter, r *http.Request) {
	Render(w, "index.html")
}

func Render(w http.ResponseWriter, tmpl string) {
	tmpl = fmt.Sprintf("public/%s", tmpl)
	t, err := template.ParseFiles(tmpl)
	if err != nil {
		log.Print("Template parsing error: ", err)
	}
	err = t.Execute(w, "")
	if err != nil {
		log.Print("Template executing error: ", err)
	}
}
