package hello

import (
	"appengine"
	"appengine/memcache"
	"fmt"
	"net/http"
)

func init() {
	http.HandleFunc("/", hello)
	http.HandleFunc("/hits", showHits)
}

func hello(w http.ResponseWriter, r *http.Request) {
	name := r.Header.Get("name")
	fmt.Fprintf(w, "Hello %s!", name)
}

func showHits(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "%d", hits(r))
}

func hits(r *http.Request) uint64 {
        c := appengine.NewContext(r)
	newValue, _ := memcache.Increment(c, "hits", 1, 0)
	return newValue
}
