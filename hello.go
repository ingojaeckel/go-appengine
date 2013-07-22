package hello

import (
	"appengine"
	"appengine/memcache"
	"fmt"
	"net/http"
)

func init() {
	http.HandleFunc("/", hello)
	http.HandleFunc("/no/content", noContent)
	http.HandleFunc("/no/content/cache", noContentCache)
	http.HandleFunc("/hits", showHits)
}

func noContent(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(204)
}

func noContentCache(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Cache-Control", "public, max-age=10")
	w.Header().Set("Pragma", "Public")
	w.WriteHeader(204)
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
