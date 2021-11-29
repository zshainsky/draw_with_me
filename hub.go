package draw

import (
	"fmt"
)

type Hub struct {
	clients map[string]*Client
}

func NewHub() *Hub {
	return &Hub{
		clients: make(map[string]*Client),
	}
}

func (h *Hub) RegisterClient(c *Client) {
	h.clients[c.id] = c
}

func (h *Hub) UnregisterClient(c *Client) {
	_, ok := h.clients[c.id]
	if ok {
		delete(h.clients, c.id)
	}
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
