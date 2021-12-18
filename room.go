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

const htmlFileName = "../room.html"

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

	room.router.Handle(fmt.Sprintf("/room-%v", room.Id), AuthMiddleware(ServeRoom))
	room.router.HandleFunc(fmt.Sprintf("/room-%v/ws", room.Id), func(w http.ResponseWriter, r *http.Request) {
		// create and serve client
		ServeWS(room.Hub, w, r)
	})
}

func ServeRoom(w http.ResponseWriter, r *http.Request) {
	log.Println(r.URL)

	fmt.Printf("/room cookies: %+v", r.Cookies())

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
