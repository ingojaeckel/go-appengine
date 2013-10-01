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

	r.HandleFunc("/rest/join", join)
	r.HandleFunc("/rest/poll", poll)
	r.HandleFunc("/rest/move/{uuid}/{x:[0-9]+}/{y:[0-9]+}", move)
	r.HandleFunc("/rest/json", getJson)
	r.HandleFunc("/rest/put", memcachePut)
	r.HandleFunc("/rest/read", memcacheRead)

	http.Handle("/", r)
}

func join(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	playerCount, _ := memcache.Increment(c, "playerCount", 1, 0)
	uuid := uuid.New()

	uuidKey := fmt.Sprintf("player%d", playerCount)
	posKey := fmt.Sprintf("%v.pos", uuid)
	c.Warningf("playerCount %d", playerCount)
	c.Warningf("UUID Key %s", uuidKey)
	c.Warningf("Pos key %s", posKey)

	memcache.JSON.Set(c, &memcache.Item {
		Key: uuidKey,
		Object: uuid,
	})
	memcache.JSON.Set(c, &memcache.Item {
		Key: posKey,
		Object: Position{100, 100},
	})

	fmt.Fprintf(w, uuid)
}

func move(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	uuid := vars["uuid"]
	x, _ := strconv.Atoi(vars["x"])
	y, _ := strconv.Atoi(vars["y"])

	c := appengine.NewContext(r)
	
	posKey := fmt.Sprintf("%v.pos", uuid)
	c.Warningf("Pos key %s", posKey)

	memcache.JSON.Set(c, &memcache.Item {
		Key: posKey,
		Object: Position{x,y},
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

func memcachePut(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	memcache.JSON.Set(c, &memcache.Item{Key: "key1", Object: "Value"})
	memcache.JSON.Set(c, &memcache.Item{Key: "key2", Object: int(23)})

	w.WriteHeader(204)
}

func memcacheRead(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)

	var value1 string
	var value2 int

	memcache.JSON.Get(c, "key1", &value1)
	memcache.JSON.Get(c, "key2", &value2)

	b1, _ := json.Marshal(value1)
	fmt.Fprintf(w, string(b1))
	b2, _ := json.Marshal(value2)
	fmt.Fprintf(w, string(b2));
}

func getCoordinates(r * http.Request, key string) Position {
	c := appengine.NewContext(r)

	var uuid string
	memcache.JSON.Get(c, key, &uuid)

	posKey := fmt.Sprintf("%v.pos", uuid)

	c.Warningf("UUID %s", uuid)
	c.Warningf("pos key %s", posKey)

	var pos Position

	memcache.JSON.Get(c, posKey, &pos)

	return pos
}

func poll(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	
	var playerCount int
	
	memcache.JSON.Get(c, "playerCount", &playerCount)

	if 1 == playerCount {
		pos := getCoordinates(r, "player1")
		s := State{1, pos, Position{0,0}}

		bytes, _ := json.Marshal(s)
		fmt.Fprintf(w, string(bytes))
	} else {
		pos1 := getCoordinates(r, "player1")
		pos2 := getCoordinates(r, "player1")

		s := State{2, pos1, pos2}
		bytes, _ := json.Marshal(s)
		fmt.Fprintf(w, string(bytes))
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

