// Runs a string existence filter server
//
// The server can be queried on the host and port specified.
//
package main

import (
	"flag"
	"io"
	"log"
	"net/http"
	"strconv"

	"github.com/wbbradley/go-exist/filter"
)

var (
	keyFilter = filter.NewMapFilter()
	host      = flag.String("host", "", "Specify the server host to listen on")
	port      = flag.Int("port", 8001, "Specify the server port to listen on")
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
	flag.Parse()

	// Set up some fake data
	filter.ReadKeysIntoFilter(keyFilter, []string{"a", "b"})

	// Set up the web handler
	http.HandleFunc("/exists", queryServer)
	err := http.ListenAndServe(*host+":"+strconv.Itoa(*port), nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
