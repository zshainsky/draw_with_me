package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/zshainsky/draw-with-me"
)

var rooms []*draw.Room

const htmlFileName = "../home.html"

func serveHome(w http.ResponseWriter, r *http.Request) {
	log.Println(r.URL)
	if r.URL.Path == "/get-rooms" {
		type roomsJSON struct {
			RoomsList []string
		}
		roomIds := []string{}
		if len(rooms) > 0 {
			for _, room := range rooms {
				roomIds = append(roomIds, room.Id)
			}
			response := roomsJSON{
				RoomsList: roomIds,
			}

			responsJSON, err := json.Marshal(response)
			if err != nil {
				fmt.Printf("get-rooms: could not create json string to return in responseText")
			}

			fmt.Printf("sending response to get-rooms: %v\n", string(responsJSON))
			w.Write([]byte(responsJSON))
		}
		return
	}
	if r.URL.Path == "/create-room" {
		// create unique room and start hub
		room := draw.NewRoom()
		rooms = append(rooms, room)

		fmt.Printf("created room in ServeHome function Handler id(room-%v)\n", room.Id)
		room.StartRoom()

		// write room id (url) back to the server
		w.Write([]byte("/room-" + room.Id))
		return
	}
	if r.URL.Path != "/" {
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	http.ServeFile(w, r, htmlFileName)
}

func main() {
	http.HandleFunc("/", serveHome)
	http.HandleFunc("/create-room", serveHome)

	// go func(rooms []*draw.Room) {
	// 	rooms = append(rooms, <-draw.roomChan)
	// }(rooms)

	log.Fatal(http.ListenAndServe(":3000", nil))

}
