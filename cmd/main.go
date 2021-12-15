package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/zshainsky/draw-with-me"
	"google.golang.org/api/idtoken"
)

var rooms []*draw.Room
var router *mux.Router

type roomsJSON struct {
	RoomsList []draw.RoomJSON
}

type AuthRequestBody struct {
	Credential string
}

type APIResponse struct {
	Code int //should be http.<response code>
}

const googleClientId = "406504108908-4djtjr6q3lil4rgrnbjproqi7ruc59vs.apps.googleusercontent.com"

const htmlFileName = "../home.html"
const htmlSigninFileName = "../signin.html"

func serveHome(w http.ResponseWriter, r *http.Request) {
	log.Println(r.URL)
	// write list of rooms as a response to the api call in json format
	// TODO: This is redundant...remove the fore loop and replace with roomsJSON{RoomsList: rooms}. Only need to pass in the reference to the globabl variable rooms
	if r.URL.Path == "/get-rooms" {
		roomIds := []draw.RoomJSON{}
		response := roomsJSON{
			RoomsList: roomIds,
		}
		if len(rooms) > 0 {
			// add all rooms to roomsIds list
			for _, room := range rooms {
				roomIds = append(roomIds, draw.RoomJSON{
					Id: room.Id,
				})
			}
			// use struct roomsJSON to format json
			response.RoomsList = roomIds
		}
		// write struct as json string
		responsJSON, err := json.Marshal(response)
		if err != nil {
			fmt.Printf("get-rooms: could not create json string to return in responseText")
		}

		fmt.Printf("sending response to get-rooms: %v\n", string(responsJSON))
		w.Write([]byte(responsJSON))

		return
	}
	if r.URL.Path == "/create-room" {

		// create unique room and start hub
		room := draw.NewRoom(router)
		rooms = append(rooms, room)

		fmt.Printf("created room in ServeHome function Handler id(room-%v)\n", room.Id)
		room.StartRoom()

		response := draw.RoomJSON{
			Id: room.Id,
		}
		// write struct as json string
		responseJSON, err := json.Marshal(response)
		if err != nil {
			fmt.Printf("get-rooms: could not create json string to return in responseText")
		}

		// write room id (url) back to the server
		w.Write([]byte(responseJSON))
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

func serveSignin(w http.ResponseWriter, r *http.Request) {
	// TODO: Break out /authorize end point into own handler  because may be called from outside of the signin flow
	if r.URL.Path == "/authorize" && r.Method == "POST" {
		reqBodyJSON := AuthRequestBody{}
		// convert request body into []byte
		err := json.NewDecoder(r.Body).Decode(&reqBodyJSON)
		if err != nil {
			panic(err)
		}
		fmt.Printf("%+v\n", reqBodyJSON.Credential)

		if reqBodyJSON.Credential == "" {
			fmt.Printf("issue getting Credential from request body\n")
		}

		payload, err := idtoken.Validate(context.Background(), string(reqBodyJSON.Credential), googleClientId)
		if err != nil {
			panic(err)
		}
		fmt.Print(payload.Claims)

		// TODO: Perform auth token verification
		// TODO: Check if user exists in database
		// TODO: if user exists --> get list of rooms for this user --> Return list to homepage
		// TODO: create user if not exists --> Respond on success --> go to homepage

		// format response to frontend as JSON object
		responseJSON, err := json.Marshal(APIResponse{
			Code: http.StatusOK,
		})
		if err != nil {
			fmt.Printf("get-rooms: could not create json string to return in responseText")
		}
		// fmt.Printf("Send JSON response:\n%s\n", responseJSON)

		// write response to fron end
		w.Write([]byte(responseJSON))

		return
	}
	if r.URL.Path != "/signin" && r.URL.Path != "/authorize" {
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	http.ServeFile(w, r, htmlSigninFileName)
}

func main() {
	router = mux.NewRouter()
	router.HandleFunc("/", serveHome)
	router.HandleFunc("/get-rooms", serveHome)
	router.HandleFunc("/create-room", serveHome)
	router.PathPrefix("/lib/").Handler(
		http.StripPrefix("/lib/", http.FileServer(http.Dir("lib/"))),
	)

	router.HandleFunc("/signin", serveSignin)
	router.HandleFunc("/authorize", serveSignin)

	// go func(rooms []*draw.Room) {
	// 	rooms = append(rooms, <-draw.roomChan)
	// }(rooms)

	err := http.ListenAndServe(":3000", router)
	if err != nil {
		log.Fatal(err)
	}

}
