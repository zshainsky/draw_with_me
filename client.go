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
}

func NewClient(conn *websocket.Conn) *Client {
	id, err := uuid.NewV4()
	if err != nil {
		fmt.Printf("problem creating unique id for client, %v", err)
		return nil
	}

	return &Client{
		id:   id.String(),
		conn: conn,
		send: make(chan []byte),
	}
}

func (c *Client) GetId() string {
	return c.id
}

// func (c *Client) wsHandler(w http.ResponseWriter, r *http.Request) {
// 	ws := newClientWS(w, r)
// 	_, msg, err := ws.ReadMessage()
// 	if err != nil {
// 		log.Printf("error reading from websocket %v\n", err)
// 	}
// 	fmt.Printf("message from websocket: ", msg)
// }

// func newClientWS(w http.ResponseWriter, r *http.Request) *websocket.Conn {
// 	conn, err := wsUpgrader.Upgrade(w, r, nil)
// 	if err != nil {
// 		log.Printf("prolem upgrading connection to WebSockets %v\n", err)
// 	}

// 	return conn

// }

func PrintClient() string {
	return "client"
}
