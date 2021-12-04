package draw

import (
	"fmt"
)

type Hub struct {
	clients    map[string]*Client
	register   chan *Client
	unregister chan *Client
	broadcast  chan []byte
}

func NewHub() *Hub {
	return &Hub{
		clients:    make(map[string]*Client),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		broadcast:  make(chan []byte),
	}
}

func (h *Hub) Run() {
	fmt.Println("Starting Hub...")
	// i := 0
	for {
		select {
		case client := <-h.register:
			h.RegisterClient(client)
		case client := <-h.unregister:
			h.UnregisterClient(client)
		case payload := <-h.broadcast:
			h.BroadcastPayload(payload)
		}

		// fmt.Printf("hub iteration: %d\n", i)
		// i += 1
	}
}

func (h *Hub) BroadcastPayload(payload []byte) {
	for _, client := range h.clients {
		select {
		case client.send <- payload:
			fmt.Printf("payload sent to client (%v): %v\n", client.GetId(), string(payload))
		default:
			fmt.Printf("default: %v\n", client.GetId())
			h.UnregisterClient(client)
		}
	}
}

func (h *Hub) RegisterClient(c *Client) error {
	if c == nil {
		return fmt.Errorf("rigister Client when Client is nil %d\n", 0)
	}
	h.clients[c.id] = c
	return nil
}

func (h *Hub) UnregisterClient(c *Client) error {
	if c == nil {
		return fmt.Errorf("unregister Client when Client is nil %d\n", 0)
	}

	_, ok := h.clients[c.id]
	if ok {
		delete(h.clients, c.id)
		close(c.send)
	}
	return nil
}

func (h *Hub) GetRegistrationChan() chan *Client {
	return h.register
}
func (h *Hub) GetUnregistrationChan() chan *Client {
	return h.unregister
}

func (h *Hub) GetClient(c *Client) (*Client, error) {
	if h.clients[c.id] == nil {
		return nil, fmt.Errorf("Hub cannot find Client with ID %v", c.id)
	}
	return h.clients[c.id], nil
}

func (h *Hub) GetNumClients() int {
	return len(h.clients)
}

func PrintHub() string {
	return "Hub"
}
