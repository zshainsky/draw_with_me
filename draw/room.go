package draw

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
	uuid "github.com/nu7hatch/gouuid"
)

type Room struct {
	Id       string
	Hub      *Hub
	HTMLFile string
	router   *mux.Router
}

type RoomJSON struct {
	Id string `json:id`
}

const htmlFileName = "static/html/room.html"

func NewRoom(r *mux.Router) *Room {
	id, err := uuid.NewV4()
	if err != nil {
		fmt.Printf("problem creating unique id for client, %v", err)
		return nil
	}

	return &Room{
		Id:       id.String(),
		Hub:      nil, // Hub is set when the room is started
		HTMLFile: htmlFileName,
		router:   r,
	}
}

func (room *Room) StartRoom() {
	// create, assign and run hub
	hub := NewHub()
	room.Hub = hub

	go hub.Run()
	room.CreateRoomRoutes()

}

func (room *Room) CreateRoomRoutes() {
	room.router.HandleFunc(fmt.Sprintf("/room-%v", room.Id), room.ServeRoom)
	// Set up route to handle websocket
	room.router.Handle(fmt.Sprintf("/room-%v/ws", room.Id), AuthMiddleware(func(w http.ResponseWriter, r *http.Request) {
		fmt.Printf("connecting to websocket...\n")
		ServeWS(room, w, r)
	}))
}

func (room *Room) ServeRoom(w http.ResponseWriter, r *http.Request) {
	log.Println(r.URL)

	if !strings.Contains(r.URL.Path, "/room") {
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	http.ServeFile(w, r, htmlFileName)
}
