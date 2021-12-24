package draw

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"google.golang.org/api/idtoken"
)

var rooms []*Room
var users map[string]*User

type Server struct {
	router *mux.Router
}

type roomsJSON struct {
	RoomsList []RoomJSON
}

const htmlHomeFileName = "static/html/home.html"
const htmlSigninFileName = "static/html/signin.html"

func NewServer(r *mux.Router) *Server {
	users = make(map[string]*User)

	server := &Server{
		router: r,
	}

	server.createRoutes()

	return server
}

func (s *Server) createRoutes() {
	s.router.Handle("/", AuthMiddleware(serveHome))

	s.router.HandleFunc("/signin", serveSignin)
	s.router.HandleFunc("/authorize", serveSignin)

	s.router.Handle("/get-rooms", AuthMiddleware(s.serveRoomActions))
	s.router.Handle("/create-room", AuthMiddleware(s.serveRoomActions))

	s.router.Handle("/user-info", AuthMiddleware(s.serveUserInfo))

	s.router.PathPrefix("/static/").Handler(
		http.StripPrefix("/static/", http.FileServer(http.Dir("static/"))), //http.Dir("") is a relative path from where the router is created. In this case it is created in main main.go in the top level dir"
	)
}

func serveHome(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("serveHome: %v", r.URL)
	// Access context values in handlers like this - Use Payload information to route request to the database for this specifci users
	payload := r.Context().Value(CTXKey("jwt"))
	_, err := getCurrentUser(w, payload)
	if err != nil {
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

	http.ServeFile(w, r, htmlHomeFileName)
}

func serveSignin(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("serveSignin: %v\n", r.URL)
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
		// Set the jwt token as a cookie upon redirect for the client
		cookie := http.Cookie{
			Name:   "jwt-token",
			Value:  "Bearer " + credential,
			Path:   "/",
			Secure: true,
		}
		http.SetCookie(w, &cookie)
		// // Redirect page after checks have been completed and user is logged in
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

func (s *Server) serveRoomActions(w http.ResponseWriter, r *http.Request) {
	payload := r.Context().Value(CTXKey("jwt"))
	targetUser, err := getCurrentUser(w, payload)
	if err != nil {
		return
	}
	// write list of rooms as a response to the api call in json format
	if r.URL.Path == "/get-rooms" {
		// fmt.Printf("/get-rooms cookies: %+v", r.Cookies())
		roomIds := []RoomJSON{}
		response := roomsJSON{
			RoomsList: roomIds,
		}
		if len(targetUser.RoomsList) > 0 {
			// add all rooms to roomsIds list
			for _, room := range targetUser.RoomsList {
				roomIds = append(roomIds, RoomJSON{
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
		w.Header().Add("Content-Type", "application/json")
		w.Write([]byte(responsJSON))

		return
	}
	if r.URL.Path == "/create-room" {
		// fmt.Printf("/create-room cookies: %+v", r.Cookies())

		// create unique room and start hub
		room := NewRoom(s.router)

		// add room to target user's room list and global rooms list
		targetUser.AddRoom(room)
		rooms = append(rooms, room)

		// fmt.Printf("created room in ServeHome function Handler id(room-%v)\n", room.Id)
		room.StartRoom()

		response := RoomJSON{
			Id: room.Id,
		}
		// write struct as json string
		responseJSON, err := json.Marshal(response)
		if err != nil {
			fmt.Printf("get-rooms: could not create json string to return in responseText")
		}
		w.Header().Add("Content-Type", "application/json")
		// write room id (url) back to the server
		w.Write([]byte(responseJSON))
		return
	}
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

func (s *Server) serveUserInfo(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("serveUserInfo: %v", r.URL)

	payload := r.Context().Value(CTXKey("jwt"))
	targetUser, err := getCurrentUser(w, payload)
	if err != nil {
		return
	}
	// Access context values in handlers like this - Use Payload information to route request to the database for this specifci users

	if r.URL.Path == "/user-info" && r.Method == "GET" {
		userJSONEvent := UserJSONEvents{
			Name:    targetUser.name,
			Email:   targetUser.email,
			Picture: targetUser.picture,
		}
		responseJSON, err := json.Marshal(userJSONEvent)
		if err != nil {
			fmt.Printf("get-rooms: could not create json string to return in responseText")
		}
		w.Header().Add("Content-Type", "application/json")
		// write room id (url) back to the server
		w.Write([]byte(responseJSON))
		return
	}
	if r.URL.Path != "/user-info" {
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
}

// Return the user from the payload token. If user not in data store, create user and return user. Return error if payload is nil, payload is not type *idtoken.Payload or if error checking on the user in data store
func getCurrentUser(w http.ResponseWriter, payload interface{}) (*User, error) {
	if payload != nil {
		var isPayload bool
		// cast payload to *idtoken.Payload
		tokenPayload, isPayload := payload.(*idtoken.Payload)
		// if cast failed
		if !isPayload {
			// fmt.Errorf("could not cast payload to *idtoken.Payload:  %v", tokenPayload)
			http.Error(w, "Not found", http.StatusNotFound)
			return nil, fmt.Errorf("could not cast payload to *idtoken.Payload:  %v", tokenPayload)
		}
		var isUserExist bool
		// check if user exists, if not, add user to in memory store
		targetUser, isUserExist, err := checkIfUserExists(tokenPayload)
		if err != nil {
			return nil, err
		}
		if !isUserExist {
			fmt.Printf("\ncreating user handler...\n")
			targetUser, _ = addUserToInMemoryStore(tokenPayload)
		}
		fmt.Printf("\npayload check complete:\n")
		fmt.Printf("%v \n\tuserName:  %v\n", targetUser.id, targetUser.email)
		return targetUser, nil
	}
	return nil, fmt.Errorf("payload is nil")
}

func checkIfUserExists(tokenPayload *idtoken.Payload) (*User, bool, error) {
	// In order to serve these routes we must have a payload from our middleware, otherwise we should never get to this point
	if tokenPayload != nil {
		// Check if user exists
		user, userExists := users[tokenPayload.Subject]
		return user, userExists, nil
	}
	return nil, false, fmt.Errorf("could not extract jwt context payload")
}
func addUserToInMemoryStore(tokenPayload *idtoken.Payload) (*User, error) {
	if tokenPayload != nil {
		name := tokenPayload.Claims["name"].(string)
		email := tokenPayload.Claims["email"].(string)
		picture := tokenPayload.Claims["picture"].(string)
		subject := tokenPayload.Subject

		targetUser := NewUser(subject, AuthType("google"), name, email, picture)
		fmt.Printf("\naddUserToInMemoryStore(): targetUser: %v\n", targetUser)
		users[subject] = targetUser
		return targetUser, nil
	}
	return nil, fmt.Errorf("could not add user. Payload is nil")
}
