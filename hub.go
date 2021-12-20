package draw

import (
	"encoding/json"
	"fmt"
)

type Hub struct {
	clients        map[string]*Client // Clients are created when a web browser has loaded the room's page represent
	register       chan *Client
	unregister     chan *Client
	broadcast      chan []byte
	canvasInMemory []*PaintEvent
}
type paintEvent struct {
	CurX   float64
	CurY   float64
	LastX  float64
	LastY  float64
	Color  string
	UserId string
}
type PaintEvent struct {
	paintEvent
}
type UserJSONEvents struct {
	Name    string `json:name`
	Email   string `json:email`
	Picture string `json:picture`
}

func NewHub() *Hub {
	return &Hub{
		clients:        make(map[string]*Client),
		register:       make(chan *Client),
		unregister:     make(chan *Client),
		broadcast:      make(chan []byte),
		canvasInMemory: []*PaintEvent{},
	}
}

func (h *Hub) Run() {
	fmt.Println("Starting Hub...")
	for {
		select {
		case client := <-h.register:
			h.RegisterClient(client)
		case client := <-h.unregister:
			h.UnregisterClient(client)
		case payload := <-h.broadcast:
			// go func(p []byte) {
			// 	time.Sleep(time.Millisecond * 500)
			// }(payload)
			h.writePaintEvent(payload)
			h.BroadcastPayload(payload)
		}
	}
}

func (h *Hub) RegisterClient(c *Client) error {
	if c == nil {
		return fmt.Errorf("rigister Client when Client is nil %d\n", 0)
	}
	// Add client to the list of active clients maintained by hub and send current canvas state
	h.clients[c.id] = c
	h.sendCanvasState(c)

	// update list of active users on webpage
	h.sendActiveUserList()

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
	// update list of active users on webpage
	h.sendActiveUserList()
	return nil
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

// In memory store
// Load payload data into a PaintEvent struct which is sent as proper JSON from the frontend and stored as a []byte
func (h *Hub) writePaintEvent(payload []byte) {
	var event = &PaintEvent{}
	if err := json.Unmarshal(payload, event); err != nil {
		panic(err)
	}
	fmt.Printf("PaintEvent: %v\n", event)

	h.canvasInMemory = append(h.canvasInMemory, event)
}

// Send from in memory store
// Send the canvas state (all paint events stored in h.canvasInMemory) to client send chan []byte
func (h *Hub) sendCanvasState(c *Client) {
	responsJSON, err := json.Marshal(h.canvasInMemory)
	if err != nil {
		fmt.Printf("get-rooms: could not create json string to return in responseText")
	}
	fmt.Printf("\nsendCanvasState(): %v\n", string(responsJSON))
	c.send <- responsJSON
}

func (h *Hub) sendActiveUserList() {
	fmt.Printf("Active User List: \n")
	key := "AciveUsers"
	userJSONList := make(map[string][]UserJSONEvents)
	// resposneJSON, err := json.Marshal(h.clients)
	for _, client := range h.clients {
		userJSONList[key] = append(userJSONList[key], UserJSONEvents{
			Name:    client.user.name,
			Email:   client.user.email,
			Picture: client.user.picture,
		})
		resposneJSON, err := json.Marshal(userJSONList)
		if err != nil {
			fmt.Printf("get-rooms: could not create json string to return in responseText")
		}
		fmt.Printf("response: %+v\n", string(resposneJSON))

	}
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
