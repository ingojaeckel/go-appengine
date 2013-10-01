package hello

import (
	"appengine"
	"encoding/json"
	"appengine/memcache"
	"fmt"
	"strings"
	"net/http"
	"strconv"
	"code.google.com/p/go-uuid/uuid"
	"github.com/gorilla/mux"
)

func init() {
	r := mux.NewRouter()

	r.HandleFunc("/rest/hello", hello)
	r.HandleFunc("/rest/join", join)
	r.HandleFunc("/rest/poll", poll)
	r.HandleFunc("/rest/move/{uuid}/{x:[0-9]+}/{y:[0-9]+}", move)
	r.HandleFunc("/rest/no/content", noContent)
	r.HandleFunc("/rest/no/content/cache", noContentCache)
	r.HandleFunc("/rest/hits", showHits)
	r.HandleFunc("/rest/state/init", initState)
	r.HandleFunc("/rest/state/get", getState)
	r.HandleFunc("/rest/json", getJson)

	http.Handle("/", r)
}

type LeaderboardEntry struct {
	player string
	value int
}

type Leaderboard struct {
	entries []LeaderboardEntry
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

func join(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	playerCount, _ := memcache.Increment(c, "playerCount", 1, 0)
	uuid := uuid.New()

	memcache.JSON.Set(c, &memcache.Item {
		Key: fmt.Sprintf("player%d", playerCount),
		Object: uuid,
	})
	memcache.JSON.Set(c, &memcache.Item {
		Key: fmt.Sprintf("%v.x", uuid),
		Object: 100,
	})
	memcache.JSON.Set(c, &memcache.Item {
		Key: fmt.Sprintf("%v.y", uuid),
		Object: 100,
	})


	c.Infof("playerCount %d", playerCount)
	fmt.Fprintf(w, uuid)
}

func move(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	uuid := vars["uuid"]
	x, _ := strconv.Atoi(vars["x"])
	y, _ := strconv.Atoi(vars["y"])

	c := appengine.NewContext(r)
	
	memcache.JSON.Set(c, &memcache.Item {
		Key: fmt.Sprintf("%v.x", uuid),
		Object: x,
	})
	memcache.JSON.Set(c, &memcache.Item {
		Key: fmt.Sprintf("%v.y", uuid),
		Object: y,
	})

	w.WriteHeader(204)
}

type Position struct {
	X,Y int
}

type State struct {
	Count int
	Player1, Player2 Position
}

func getJson(w http.ResponseWriter, r *http.Request) {
	s := State{23, Position{1,2}, Position{3,4}}

	bytes, _ := json.Marshal(s)
	fmt.Fprint(w, string(bytes))
}

func getCoordinates(r * http.Request, key string) (int, int) {
	c := appengine.NewContext(r)

	var uuid string
	memcache.JSON.Get(c, key, &uuid)

	keyX1 := fmt.Sprintf("%v.x", uuid)
	keyY1 := fmt.Sprintf("%v.y", uuid)

	var x,y int

	memcache.JSON.Get(c, keyX1, &x)
	memcache.JSON.Get(c, keyY1, &y)

	return x,y
}

func poll(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	
	var playerCount int
	
	memcache.JSON.Get(c, "playerCount", &playerCount)

	if 1 == playerCount {
		x, y := getCoordinates(r, "player1")
		s := State{1, Position{x,y}, Position{}}

		// fmt.Printf(s)
		fmt.Fprintf(w, "%q", s)
	} else {
		x1, y1 := getCoordinates(r, "player1")
		x2, y2 := getCoordinates(r, "player1")

		s := State{1, Position{x1, y1}, Position{x2, y2}}

		// fmt.Printf(s)
		fmt.Fprintf(w, "%+v", s)
	}
}

func getPosition(value string) Position {
	fmt.Printf("value = %s", value)
	parts := strings.Split(value, " ")
	fmt.Printf("len %d", len(parts))
	x, _ := strconv.Atoi(parts[0])
	y, _ := strconv.Atoi(parts[1])

	return Position{x, y}
}

func showHits(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "%d", hits(r))
}

func hits(r *http.Request) uint64 {
        c := appengine.NewContext(r)
	newValue, _ := memcache.Increment(c, "hits", 1, 0)
	return newValue
}
