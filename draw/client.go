package draw

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
	uuid "github.com/nu7hatch/gouuid"
	"google.golang.org/api/idtoken"
)

var wsUpgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

type Client struct {
	id   string
	conn *websocket.Conn
	send chan []byte
	hub  *Hub
	user *User
	// room *Room // used to associate paint events
}

func NewClient(h *Hub, user *User, conn *websocket.Conn) *Client {
	id, err := uuid.NewV4()
	if err != nil {
		fmt.Printf("problem creating unique id for client, %v", err)
		return nil
	}

	return &Client{
		id:   id.String(),
		conn: conn,
		send: make(chan []byte, 256),
		hub:  h,
		user: user,
		// room: room,
	}
}

func ServeWS(room *Room, w http.ResponseWriter, r *http.Request) {

	var targetUser *User
	payload := r.Context().Value(CTXKey("jwt"))

	if payload != nil {
		var isPayload bool
		// cast payload to *idtoken.Payload
		tokenPayload, isPayload := payload.(*idtoken.Payload)
		// if cast failed
		if !isPayload {
			fmt.Errorf("could not cast payload to *idtoken.Payload:  %v", tokenPayload)
			http.Error(w, "Not found", http.StatusNotFound)
			return
		}
		var isUserExist bool
		// check if user exists, if not, add user to in memory store
		targetUser, isUserExist, _ = checkIfUserExists(tokenPayload)
		if !isUserExist {
			fmt.Printf("\ncreating user in the ServeRoom handler...\n")
			targetUser, _ = addUserToInMemoryStore(tokenPayload)
		}
		fmt.Printf("\npayload check complete:\n")
		fmt.Printf("%v \n\tuserName:  %v\n", targetUser.id, targetUser.email)
	}

	// If user visits a room created by a different user, add room to user's map of rooms
	targetUser.AddRoom(room)

	// fmt.Printf("\n/room- cookies: %+v\n", r.Cookies())
	fmt.Printf("\n/Setting up Websocket for user: %v (%v)\n", targetUser.email, targetUser.id)
	// Create websocket connection
	conn, err := wsUpgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("prolem upgrading connection to WebSockets %v\n", err)
		return
	}

	// Create new Client and register with hub
	// client := NewClient(room.Hub, targetUser, room, conn)
	client := NewClient(room.Hub, targetUser, conn)
	fmt.Printf("websocket new client: %v", client.user.email)
	room.Hub.register <- client

	go client.sendToHub()
	go client.writeToWS()
}

// sendToHub() writes messages from the websocket to the hub
func (c *Client) sendToHub() {
	// clean up connections on complete
	defer func() {
		c.hub.unregister <- c
		c.conn.Close()
	}()

	// Read message from the websocket and send to the hub for broadcasting, if issue reading from hub, break and close connections
	for {
		_, msg, err := c.conn.ReadMessage()
		if err != nil {
			log.Printf("could not read message from ws: %v", err)
			break
		}
		// print message read from websocket
		// fmt.Printf("websocket message: %s\n", msg)

		//TODO: Add the RoomID to the message recived from the Websocket. This is likely set in lib/components/room-canvas.js
		c.hub.broadcast <- msg
	}
}

func (c *Client) writeToWS() {
	// close connection to ws after completion
	defer func() {
		c.conn.Close()
	}()
	// run infinate loop waiting (blocking) for messages from the hub (<-c.send)
	for {
		payload := <-c.send
		// fmt.Printf("client send chan: %v\n", string(payload))
		// write payload from hub to ws, if there is an error break out of loop and close connection
		err := c.conn.WriteMessage(1, payload)
		if err != nil {
			log.Printf("could not write message to ws: %v\n", err)
			break
		}
	}
}

func (c *Client) SendUpdate(payload string) {
	c.hub.broadcast <- []byte(payload)
}

func (c *Client) SendRegistration() {
	c.hub.register <- c
}

func (c *Client) SendUnregistration() {
	c.hub.unregister <- c
}

func (c *Client) GetRegChan() chan *Client {
	return c.hub.register
}
func (c *Client) GetUnregChan() chan *Client {
	return c.hub.unregister
}
func (c *Client) GetSendChan() chan []byte {
	return c.send
}

func (c *Client) GetId() string {
	return c.id
}

func PrintClient() string {
	return "client"
}
