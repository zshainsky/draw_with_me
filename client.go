package draw

import (
	"fmt"

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
		conn: conn,
		send: make(chan []byte, 256),
		hub:  h,
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
