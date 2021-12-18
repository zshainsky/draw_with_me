package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/zshainsky/draw-with-me"
)

var rooms []*draw.Room
var router *mux.Router

type roomsJSON struct {
	RoomsList []draw.RoomJSON
}

type AuthRequestBody struct {
	Credential   string `json:credential`
	G_csrf_token string `json:g_csrf_token`
}

type APIResponse struct {
	Code int //should be http.<response code>
}

const googleClientId = "406504108908-4djtjr6q3lil4rgrnbjproqi7ruc59vs.apps.googleusercontent.com"

const htmlFileName = "../home.html"
const htmlSigninFileName = "../signin.html"

func serveHome(w http.ResponseWriter, r *http.Request) {
	log.Printf("serveHome: %v", r.URL)
	// Access context values in handlers like this - Use Payload information to route request to the database for this specifci users
	payload := r.Context().Value("jwt_payload")
	fmt.Printf("\n serveHome: Contexts: %+v \n", payload)
	// write list of rooms as a response to the api call in json format
	// TODO: This is redundant...remove the fore loop and replace with roomsJSON{RoomsList: rooms}. Only need to pass in the reference to the globabl variable rooms
	if r.URL.Path == "/get-rooms" {
		fmt.Printf("/get-rooms cookies: %+v", r.Cookies())
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
		fmt.Printf("/create-room cookies: %+v", r.Cookies())

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
	log.Printf("serveSignin: %v", r.URL)
	// TODO: Break out /authorize end point into own handler  because may be called from outside of the signin flow
	if r.URL.Path == "/authorize" && r.Method == "POST" {
		// Authorization is using the data-login_uri="/authorize" tag with the google button to post this response to the webserver after user has logged in
		headerContentTtype := r.Header.Get("Content-Type")
		if headerContentTtype != "application/x-www-form-urlencoded" {
			http.Error(w, "Incompatible Content Type", http.StatusUnsupportedMediaType)
			return
		}
		checkDoubleSubmitCookie(w, r)

		// Check if credential is in Post body
		if err, ok := r.PostForm["credential"]; !ok {
			http.Error(w, fmt.Sprintf("Error parsing credential form value on auth request: %v", err), http.StatusBadRequest)
			return
		}
		credential := r.PostFormValue("credential")

		// // TODO: Check if user exists in database
		// // TODO: if user exists --> get list of rooms for this user --> Return list to homepage
		// // TODO: create user if not exists --> Respond on success --> go to homepage

		// Set the jwt token as a cookie upon redirect for the client
		cookie := http.Cookie{
			Name:   "jwt-token",
			Value:  "Bearer " + credential,
			Path:   "/",
			Secure: true,
		}
		http.SetCookie(w, &cookie)
		// Redirect page after checks have been completed and user is logged in
		http.Redirect(w, r, "/", http.StatusMovedPermanently)

		return
	}
	// redirect /authorize page to /signin page
	if r.URL.Path == "/authorize" && r.Method == "GET" {
		fmt.Println("GET /authorize")
		// Redirect page after checks have been completed and user is logged in
		http.Redirect(w, r, "/signin", http.StatusMovedPermanently)

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

func checkDoubleSubmitCookie(w http.ResponseWriter, r *http.Request) {
	// START: Verify double submit cookie:
	// 1. Check if g_csrf_token in body
	err := r.ParseForm()
	if err != nil {
		http.Error(w, fmt.Sprintf("Error parsing form values on auth request: %v", err), http.StatusBadRequest)
		return
	}
	if err, ok := r.PostForm["g_csrf_token"]; !ok {
		http.Error(w, fmt.Sprintf("No CSRF token in Body: %v", err), http.StatusBadRequest)
		return
	}
	g_csrf_body := r.PostFormValue("g_csrf_token")

	// 2. Check if in cookie:
	g_csrf_cookie, err := r.Cookie("g_csrf_token")
	if err != nil {
		http.Error(w, "No CSRF token in Cookie.", http.StatusBadRequest)
		return
	}
	fmt.Printf("Compare g_csrf_token in Cookie and Body:  %+v == %+v \n", g_csrf_cookie.Value, g_csrf_body)

	// 3. Check if cookie and body value match
	if g_csrf_cookie.Value != g_csrf_body {
		http.Error(w, "Failed to verify double submit cookie.", http.StatusBadRequest)
		return
	}
	// END: Verify double submit cookie:
}

func main() {
	router = mux.NewRouter()
	router.Handle("/", draw.AuthMiddleware(serveHome))
	router.Handle("/get-rooms", draw.AuthMiddleware(serveHome))
	router.Handle("/create-room", draw.AuthMiddleware(serveHome))

	router.HandleFunc("/signin", serveSignin)
	router.HandleFunc("/authorize", serveSignin)

	router.PathPrefix("/lib/").Handler(
		http.StripPrefix("/lib/", http.FileServer(http.Dir("lib/"))),
	)

	// go func(rooms []*draw.Room) {
	// 	rooms = append(rooms, <-draw.roomChan)
	// }(rooms)

	err := http.ListenAndServe(":3000", router)
	if err != nil {
		log.Fatal(err)
	}

}
