package main

import (
	"io"
	"log"
	"net/http"

	"github.com/wbbradley/go-exist/filter"
)

var (
	keyFilter = filter.NewMapFilter([]string{"a", "b"})
)

func queryServer(w http.ResponseWriter, req *http.Request) {
	if req.Method == "GET" {
		queryServerGet(w, req)
	} else {
		w.WriteHeader(400)
	}
}

// Handle existence queries
func queryServerGet(w http.ResponseWriter, req *http.Request) {
	query := req.FormValue("q")

	w.Header().Set("Content-Type", "application/json")

	if len(query) > 0 && keyFilter.KeyExists(query) {
		io.WriteString(w, "{\"exists\": true}")
	} else {
		io.WriteString(w, "{\"exists\": false}")
	}
}

func main() {
	http.HandleFunc("/", queryServer)
	err := http.ListenAndServe(":8001", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
