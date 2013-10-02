package hello

import (
	"appengine"
	"encoding/json"
	"appengine/memcache"
	"fmt"
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

type State struct {
	Players []Player 
}

type Player struct {
	ID, Name string
	P Position
}

type Position struct {
	X,Y int
}

func poll(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	
	var playerCount int
	
	memcache.JSON.Get(c, "playerCount", &playerCount)

	var players []Player

	switch playerCount {
		case 0:
			players = []Player{}
		case 1:
			players = []Player{getPlayer(r, "player1")}
		default:
			players = make([]Player, playerCount)
			for i := 0; i < playerCount; i++ {
				key := fmt.Sprintf("player%d", i+1)
				c.Warningf("key = %s", key)
				players[i] = getPlayer(r, key)
			}
	}

	b, _ := json.Marshal(players)
	c.Warningf("players %s", string(b))

	s := State{players}
	bytes, _ := json.Marshal(s)
	fmt.Fprintf(w, string(bytes))
}

func getPlayer(r * http.Request, key string) Player {
	c := appengine.NewContext(r)

	var uuid string
	memcache.JSON.Get(c, key, &uuid)

	posKey := fmt.Sprintf("%v.pos", uuid)

	c.Warningf("UUID %s", uuid)
	c.Warningf("pos key %s", posKey)

	var pos Position

	memcache.JSON.Get(c, posKey, &pos)

	return Player{uuid, key, pos}
}


