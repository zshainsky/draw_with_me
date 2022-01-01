package draw

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/Pallinder/go-randomdata"
	"github.com/gorilla/mux"
	uuid "github.com/nu7hatch/gouuid"
	"github.com/zshainsky/draw-with-me/db"
)

type Room struct {
	Id   string
	Name string
	Hub  *Hub
	// HTMLFile string
	router    *mux.Router
	isStarted bool
}

type RoomJSON struct {
	Id          string                   `json:id`
	Name        string                   `json:name`
	CanvasState map[string][]*PaintEvent `json:",omitempty"`
}

const htmlFileName = "static/html/room.html"

func NewRoom(r *mux.Router) *Room {
	id, err := uuid.NewV4()
	if err != nil {
		fmt.Printf("problem creating unique id for client, %v", err)
		return nil
	}

	roomName := generateRoomName()

	// DB: Add to Database: 'room'
	db.InsertRoom(db.RoomTable{
		Id:   id.String(),
		Name: roomName,
	})

	return &Room{
		Id:   id.String(),
		Name: roomName,
		Hub:  nil, // Hub is set when the room is started
		// HTMLFile: htmlFileName,
		router:    r,
		isStarted: false,
	}

}

func (room *Room) StartRoom() {
	// create, assign and run hub
	hub := NewHub(RoomJSON{
		Id:          room.Id,
		Name:        room.Name,
		CanvasState: make(map[string][]*PaintEvent),
	})
	room.Hub = hub

	go hub.Run()
	room.CreateRoomRoutes()
	room.isStarted = true

}

func (room *Room) CreateRoomRoutes() {
	room.router.Handle(fmt.Sprintf("/room-%v", room.Id), AuthMiddleware(room.ServeRoom))
	// Set up route to handle websocket
	room.router.Handle(fmt.Sprintf("/room-%v/ws", room.Id), AuthMiddleware(func(w http.ResponseWriter, r *http.Request) {
		fmt.Printf("connecting to websocket...\n")
		ServeWS(room, w, r)
	}))
}

func (room *Room) GetCurrentRoomCanvasStateJSON() []byte {
	responseJSON, err := room.Hub.GetCanvasStateAsJSON(room.Hub.canvasInMemory)
	if err != nil {
		fmt.Printf("get-rooms: could not create json string to return in responseText: %v\n", err)
		return nil
	}

	fmt.Printf("GetCurrentRoomCanvasStateJSON()\n")
	return responseJSON
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

func generateRoomName() string {
	return fmt.Sprintf("%v %v %v", randomdata.Title(randomdata.RandomGender), randomdata.Country(randomdata.FullCountry), randomdata.SillyName())
}
