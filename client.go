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

func NewClient(h *Hub, conn *websocket.Conn) *Client {
	id, err := uuid.NewV4()
	if err != nil {
		fmt.Printf("problem creating unique id for client, %v", err)
		return nil
	}

	return &Client{
		id:   id.String(),
		conn: conn,
		send: make(chan []byte, 5),
		hub:  h,
	}
}

func ServeWS(hub *Hub, w http.ResponseWriter, r *http.Request) {
	// Create websocket connection
	conn, err := wsUpgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("prolem upgrading connection to WebSockets %v\n", err)
		return
	}

	// Create new Client and register with hub
	client := NewClient(hub, conn)
	hub.register <- client

	go client.sendToHub()
	go client.writeToWS()
}

// sendToHub() writes messages from the websocket to the hub
func (c *Client) sendToHub() {
	// clean up connections
	defer func() {
		c.hub.unregister <- c
		c.conn.Close()
	}()

	// Read message from the websocket
	for {
		_, msg, err := c.conn.ReadMessage()
		if err != nil {
			log.Printf("could not read message from ws: %v", err)
		}
		fmt.Printf("websocket message: %s\n", msg)

		c.hub.broadcast <- msg
	}
}

func (c *Client) writeToWS() {
	defer func() {

	}()
	for {
		payload := <-c.send
		fmt.Printf("client send chan: %v\n", string(payload))
		err := c.conn.WriteMessage(1, payload)
		if err != nil {
			log.Printf("could not write message to ws: %v\n", err)
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
