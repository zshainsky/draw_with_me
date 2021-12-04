package draw

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
	uuid "github.com/nu7hatch/gouuid"
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
}

func NewClient(conn *websocket.Conn, h *Hub) *Client {
	id, err := uuid.NewV4()
	if err != nil {
		fmt.Printf("problem creating unique id for client, %v", err)
		return nil
	}

	return &Client{
		id:   id.String(),
		conn: nil,
		send: make(chan []byte, 256),
		hub:  h,
	}
}

func (c *Client) ServeWS(hub *Hub, w http.ResponseWriter, r *http.Request) {
	// Create websocket connection
	conn, err := wsUpgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("prolem upgrading connection to WebSockets %v\n", err)
	}

	// Read message from the user interface
	_, msg, err := conn.ReadMessage()
	if err != nil {
		log.Printf("could not read message from ws: %v", err)
	}
	fmt.Printf("websocket message: %s\n", msg)

	// Write a message to the user interface
	err = conn.WriteMessage(1, []byte("New Client"))
	if err != nil {
		log.Printf("could not write message to ws: %v", err)
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
