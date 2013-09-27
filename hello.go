package hello

import (
	"appengine"
	"appengine/memcache"
	"fmt"
	"net/http"
)

func init() {
	http.HandleFunc("/hello", hello)
	http.HandleFunc("/no/content", noContent)
	http.HandleFunc("/no/content/cache", noContentCache)
	http.HandleFunc("/hits", showHits)
	http.HandleFunc("/state/init", initState)
	http.HandleFunc("/state/get", getState)
}

type LeaderboardEntry struct {
	player string
	value int
}

type Leaderboard struct {
	entries []LeaderboardEntry
}

type State struct {
	a,b,c,d int
	state string
}

func initState(w http.ResponseWriter, r *http.Request) {
        c := appengine.NewContext(r)

	var in struct {I int;}
	in.I = 23
	item := &memcache.Item {
		Key: "s",
		Object: in,
	}
	memcache.JSON.Set(c, item)
	
	w.WriteHeader(204)
}

func getState(w http.ResponseWriter, r *http.Request) {
        c := appengine.NewContext(r)

	var out struct {I int;}
	memcache.JSON.Get(c, "s", &out)
	fmt.Fprint(w, out)
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
