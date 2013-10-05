package hello

import (
	"appengine"
	"appengine/memcache"
	"appengine/channel"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"code.google.com/p/go-uuid/uuid"
	"github.com/gorilla/mux"
)

func init() {
	r := mux.NewRouter()

	r.HandleFunc("/rest/create", create)
	r.HandleFunc("/rest/send", send)
	r.HandleFunc("/rest/join", join)
	r.HandleFunc("/rest/poll", poll)
	r.HandleFunc("/rest/move/{uuid}/{x:[0-9]+}/{y:[0-9]+}", move)

	http.Handle("/", r)
}

func getChannelToken(c appengine.Context, uuid string) string {
	token, _ := channel.Create(c, uuid)
	c.Warningf("createChannel(%s) -> %s", uuid, token)
	return token
}

func create(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	token, err := channel.Create(c, "player")
	if err != nil {
		c.Warningf("e = %s", err)
	}
	c.Warningf("createChannel(%s) -> %s", "player", token)
	fmt.Fprintf(w, token)
}

func send(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	err := channel.Send(c, "player", "message")
	if err != nil {
		c.Warningf("e = %s", err)
	}
	w.WriteHeader(204)
}

func join(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	playerCount, _ := memcache.Increment(c, "playerCount", 1, 0)
	uuid := uuid.New()

	uuidKey := fmt.Sprintf("player%d", playerCount)
	posKey := fmt.Sprintf("%v.pos", uuid)

	memcache.JSON.Set(c, &memcache.Item {
		Key: uuidKey,
		Object: uuid,
	})
	memcache.JSON.Set(c, &memcache.Item {
		Key: posKey,
		Object: Position{100, 100},
	})

	bytes, _ := json.Marshal(JoinResponse{uuid, getChannelToken(c, uuidKey)})
	fmt.Fprintf(w, string(bytes))
}

func move(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	uuid := vars["uuid"]
	x, _ := strconv.Atoi(vars["x"])
	y, _ := strconv.Atoi(vars["y"])

	c := appengine.NewContext(r)
	
	posKey := fmt.Sprintf("%v.pos", uuid)
	memcache.JSON.Set(c, &memcache.Item {
		Key: posKey,
		Object: Position{x,y},
	})
	updatedPlayer := Player{uuid, "?", Position{x,y}}

	var playerCount int
	memcache.JSON.Get(c, "playerCount", &playerCount)

	for i := 0; i < playerCount; i++ {
		key := fmt.Sprintf("player%d", i+1)
		e := channel.SendJSON(c, key, updatedPlayer)
		if e != nil {
			c.Warningf("Error while sending channel message %s", e)
		}
	}

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

type JoinResponse struct {
	UUID, ChannelToken string
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


