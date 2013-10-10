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

type NotifyRequest struct {
	ID string
	X,Y int
	Recipients []string
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
	Players []Player
}

type PlayerList struct {
	Players []Player
}

func init() {
	r := mux.NewRouter()

	r.HandleFunc("/rest/players/add/{id}/{name}/", playersAdd)
	r.HandleFunc("/rest/players/del/{id}/", playersDelete)
	r.HandleFunc("/rest/players", playersGet)
	r.HandleFunc("/rest/create", create)
	r.HandleFunc("/rest/send", send)
	r.HandleFunc("/rest/join", join)
	r.HandleFunc("/rest/poll", poll)
	r.HandleFunc("/_ah/channel/connected/", connected)
	r.HandleFunc("/_ah/channel/disconnected/", disconnected)
	r.HandleFunc("/rest/move/{uuid}/{x:[0-9]+}/{y:[0-9]+}/", move)
	r.HandleFunc("/rest/notify", notify)

	http.Handle("/", r)
}

func getChannelToken(c appengine.Context, uuid string) string {
	token, _ := channel.Create(c, uuid)
	c.Warningf("createChannel(%s) -> %s", uuid, token)
	return token
}

func playersAdd(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)

	vars := mux.Vars(r)
	id := vars["id"]
	name := vars["name"]

	var playerList PlayerList
	memcache.JSON.Get(c, "players", &playerList)
	playerExists := false
	
	newList := make([]Player, len(playerList.Players)+1)
	for i:=0; i<len(playerList.Players); i++ {
		if playerList.Players[i].ID == id {
			// This player is already in the list.
			// Don't modify the list at all to avoid re-adding the user.
			playerExists = true
			break
		}
		newList[i] = playerList.Players[i];
	}

	if playerExists {
		// just send current list to the client
		bytes, _ := json.Marshal(playerList)
		fmt.Fprintf(w, string(bytes))
	} else {
		newList[len(playerList.Players)] = Player{id,name,Position{0,0}}

		memcache.JSON.Set(c, &memcache.Item {
			Key: "players",
			Object: PlayerList{newList},
		});

		bytes, _ := json.Marshal(PlayerList{newList})
		fmt.Fprintf(w, string(bytes))
	}
}

func playersDelete(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	vars := mux.Vars(r)
	id := vars["id"]

	var playerList PlayerList
	_, e := memcache.JSON.Get(c, "players", &playerList)
	response := playerList

	if e != nil || len(playerList.Players) == 0 {
		// reply with empty list	
		response = PlayerList{[]Player{}}
	}

	index := indexOf(playerList.Players, id)
	playerExists := index >= 0

	if playerExists {
		newList := make([]Player, len(playerList.Players) - 1)

		// create a new players list but skip the player that should be deleted
		k := 0
		for i:=0; i<len(playerList.Players); i++ {
			if i == index {
				continue
			}
			newList[k] = playerList.Players[i]
			k++
		}
		response = PlayerList{newList}

		// Update memcache
		memcache.JSON.Set(c, &memcache.Item{
			Key: "players",
			Object: response,
		});
	}

	bytes, _ := json.Marshal(response)
	fmt.Fprintf(w, string(bytes))
}

// Returns the position of the player with the given ID in the players array. Or -1 if no player with this ID was found.
func indexOf(players []Player, id string) int {
	for i:=0; i<len(players); i++ {
		if players[i].ID == id {
			return i
		}
	}
	return -1
}

func playersGet(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)

	var playerList PlayerList
	_, b := memcache.JSON.Get(c, "players", &playerList)
	if b != nil {
		c.Warningf("b = %s", b)
		playerList = PlayerList{[]Player{}}
	}
	bytes, _ := json.Marshal(playerList)
	fmt.Fprintf(w, string(bytes))
}

func connected(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	clientId := r.FormValue("from")

	var playerCount int
	memcache.JSON.Get(c, "playerCount", &playerCount)
	// no need to increase the player count - it was already increased in /join 
	newPlayer := fmt.Sprintf("[0, \"%s\"]", clientId) // 0 == add player

	for i:=0; i<playerCount; i++ {
		key := fmt.Sprintf("player%d", i+1)
		channel.Send(c, key, newPlayer)
	}

	w.WriteHeader(204)
}

func disconnected(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)
	clientId := r.FormValue("from")
	
	var playerCount int
	memcache.JSON.Get(c, "playerCount", &playerCount)
	// TODO decrease player count after disconnect
	newPlayer := fmt.Sprintf("[1, \"%s\"]", clientId) // 1 == remove player

	for i:=0; i<playerCount; i++ {
		key := fmt.Sprintf("player%d", i+1)
		channel.Send(c, key, newPlayer)
	}

	w.WriteHeader(204)
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

	posKey := fmt.Sprintf("%v.pos", uuid)
	uuidKey := fmt.Sprintf("player%d", playerCount)

	memcache.JSON.Set(c, &memcache.Item {
		Key: uuidKey,
		Object: uuid,
	})
	memcache.JSON.Set(c, &memcache.Item {
		Key: posKey,
		Object: Position{100, 100},
	})

	players := make([]Player, playerCount - 1)

	for i := 0; i < int(playerCount) - 1; i++ {
		key := fmt.Sprintf("player%d", i+1)
		players[i] = getPlayer(r, key)
	}

	bytes, _ := json.Marshal(JoinResponse{uuid, getChannelToken(c, uuidKey), players})
	fmt.Fprintf(w, string(bytes))
}

func move(w http.ResponseWriter, r *http.Request) {
	c := appengine.NewContext(r)

	vars := mux.Vars(r)
	uuid := vars["uuid"]
	x, _ := strconv.Atoi(vars["x"])
	y, _ := strconv.Atoi(vars["y"])
	
	posKey := fmt.Sprintf("%v.pos", uuid)
	memcache.JSON.Set(c, &memcache.Item {
		Key: posKey,
		Object: Position{x,y},
	})
	updatedPlayer := fmt.Sprintf("[\"%s\",%d,%d]", uuid, x, y)

	var playerCount int
	memcache.JSON.Get(c, "playerCount", &playerCount)

	for i := 0; i < playerCount; i++ {
		key := fmt.Sprintf("player%d", i+1)
		channel.Send(c, key, updatedPlayer)
	}

	w.WriteHeader(204)
}

func parseNotifyRequest(r *http.Request) NotifyRequest {
	decoder := json.NewDecoder(r.Body)
	var notifyRequest NotifyRequest

	if e := decoder.Decode(&notifyRequest); e != nil {
		c := appengine.NewContext(r)
		c.Warningf("Error decoding request %s", e)
	}

	return notifyRequest
}

func notify(w http.ResponseWriter, r *http.Request) {
	notifyRequest := parseNotifyRequest(r)
	updatedPlayer := fmt.Sprintf("[\"%s\",%d,%d]", notifyRequest.ID, notifyRequest.X, notifyRequest.Y)
	c := appengine.NewContext(r)

	for i := 0; i < len(notifyRequest.Recipients); i++ {
		c.Warningf("send('%s', '%s') [%d/%d]", notifyRequest.Recipients[i], updatedPlayer, i, len(notifyRequest.Recipients))
		channel.Send(c, notifyRequest.Recipients[i], updatedPlayer)
	}

	w.WriteHeader(204)
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

	var pos Position

	memcache.JSON.Get(c, posKey, &pos)

	return Player{uuid, key, pos}
}


