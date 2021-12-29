package draw

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/zshainsky/draw-with-me/db"
)

type Hub struct {
	roomId         string
	clients        map[string]*Client // Clients are created when a web browser has loaded the room's page represent
	register       chan *Client
	unregister     chan *Client
	broadcast      chan []byte
	canvasInMemory []*PaintEvent
	autoSave       AutoSave
}

type AutoSave struct {
	duration          time.Duration // in Seconds
	ticker            *time.Ticker
	done              chan bool
	event             chan bool
	eventTriggerCount int // num events before saving
}

type PaintEvent struct {
	EvtTime int64 // Unix() = epoch time
	UserId  string
	RoomId  string
	CurX    int
	CurY    int
	LastX   int
	LastY   int
	Color   string
}

// JSON Keys
const (
	paintEventKey  string = "PaintEvent"
	canvasStateKey string = "CanvasState"
)

func NewHub(roomId string) *Hub {
	autoSaveDuration := 10    // seconds
	autoSaveEventLimit := 100 // num events

	hub := &Hub{
		roomId:         roomId,
		clients:        make(map[string]*Client),
		register:       make(chan *Client),
		unregister:     make(chan *Client),
		broadcast:      make(chan []byte),
		canvasInMemory: []*PaintEvent{},
		autoSave: AutoSave{
			duration:          time.Duration(autoSaveDuration) * time.Second,
			ticker:            time.NewTicker(time.Duration(autoSaveDuration) * time.Second),
			done:              make(chan bool),
			event:             make(chan bool),
			eventTriggerCount: autoSaveEventLimit,
		},
	}

	// DB: Create canvasInMemory from db
	hub.initCanvasFromCanvasStateTable()

	return hub
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
			h.writePaintEventInMemory(payload)
			h.BroadcastPayload(payload)
		}
	}
}

func (h *Hub) RegisterClient(c *Client) error {
	if c == nil {
		return fmt.Errorf("rigister Client when Client is nil %d\n", 0)
	}
	// check if any clients are in the room, if none, start autosave go routine
	if len(h.clients) == 0 {
		h.startAutoSaveRoutine()
	}
	// Add client to the list of active clients maintained by hub and send current canvas state
	h.clients[c.id] = c
	h.sendCanvasState(c)
	// update list of active users on webpage
	h.BroadcastPayload(h.sendActiveUserList())

	return nil
}

func (h *Hub) UnregisterClient(c *Client) error {
	if c == nil {
		return fmt.Errorf("unregister Client when Client is nil %d\n", 0)
	}

	// Remove client from clients map
	_, ok := h.clients[c.id]
	if ok {
		delete(h.clients, c.id)
		close(c.send)
	}
	// check if any clients are in the room, if none, stop autosave go routine
	if len(h.clients) == 0 {
		h.stopAutoSaveRoutine()
	}

	// when a client leaves the room, send an update to the canvas_state table in the db
	h.writeCanvasStateToDB()

	// update list of active users on webpage
	h.BroadcastPayload(h.sendActiveUserList())
	return nil
}

func (h *Hub) BroadcastPayload(payload []byte) {
	for _, client := range h.clients {
		select {
		case client.send <- payload:
			// print payload being sent to client
			// fmt.Printf("payload sent to client (%v): %v\n", client.GetId(), string(payload))
		default:
			fmt.Printf("default: %v\n", client.GetId())
			h.UnregisterClient(client)
		}
	}
}

// In memory store
// Load payload data into a PaintEvent struct which is sent as proper JSON from the frontend and stored as a []byte
func (h *Hub) writePaintEventInMemory(payload []byte) {
	paintEventMap := make(map[string]*PaintEvent)
	// paintEventMap[key] = &PaintEvent{}

	if err := json.Unmarshal(payload, &paintEventMap); err != nil {
		panic(err)
	}
	// fmt.Printf("PaintEvent: +%v\n", paintEventMap[paintEventKey])

	h.canvasInMemory = append(h.canvasInMemory, paintEventMap[paintEventKey])
	h.autoSave.event <- true
}

func (h *Hub) writeCanvasStateToDB() {
	responseJSON, err := getCanvasStateAsJSON(h.canvasInMemory)
	if err != nil {
		fmt.Printf("get-rooms: could not create json string to return in responseText")
	}
	// fmt.Printf("\n")
	// fmt.Printf("%s\n", responseJSON)
	// Create a blank canvas state map
	rowsAffected := db.UpdateCanvasStateForRoom(db.CanvasStateTable{
		RoomId:     h.roomId,
		CanvasJSON: string(responseJSON),
	})

	// No rows affected due to no record for room_id in canvas_json
	if rowsAffected == 0 {
		db.InsertCanvasStateForRoom(db.CanvasStateTable{
			RoomId:     h.roomId,
			CanvasJSON: string(responseJSON),
		})
	}

}

// Send from in memory store
// Send the canvas state (all paint events stored in h.canvasInMemory) to client send chan []byte
func (h *Hub) sendCanvasState(c *Client) {
	// userJSONList := make(map[string][]*PaintEvent)
	// userJSONList[canvasStateKey] = h.canvasInMemory

	responseJSON, err := getCanvasStateAsJSON(h.canvasInMemory)
	if err != nil {
		fmt.Printf("get-rooms: could not create json string to return in responseText: %v\n", err)
	}

	fmt.Printf("\nsendCanvasState()\n")
	// fmt.Printf("%s\n", responseJSON)
	c.send <- responseJSON
}

// Initialize canvas state in memory. This function queries the canvas_state table and loads the canvas_json value into hub.canvasInMemory.
// If there is no record for this room in the canvas_state table, then insert a blank state into canvas_json.
func (h *Hub) initCanvasFromCanvasStateTable() {
	dbCanvasState, err := db.GetCanvasStateForRoom(h.roomId)
	if err != nil {
		fmt.Printf("error initializing canvas for room (%v): %v\n", h.roomId, err)
		return
	}

	// DB: If no canvas events exist for this room, insert a blank canvas state. Each room should only have 1 record of canvas state
	if dbCanvasState == (db.CanvasStateTable{}) {
		fmt.Printf("No record for db canvas state so insert a blank record...\n")

		responseJSON, err := getCanvasStateAsJSON([]*PaintEvent{})
		if err != nil {
			fmt.Printf("get-rooms: could not create json string to return in responseText")
		}
		// fmt.Printf("\n")
		// fmt.Printf("%s\n", responseJSON)

		db.InsertCanvasStateForRoom(db.CanvasStateTable{
			RoomId:     h.roomId,
			CanvasJSON: string(responseJSON),
		})
		return
	}

	canvasStateMap := make(map[string][]*PaintEvent)
	canvasStateMap[canvasStateKey] = []*PaintEvent{}
	// Extract all paint events from CanvasState json
	if err := json.Unmarshal([]byte(dbCanvasState.CanvasJSON), &canvasStateMap); err != nil {
		fmt.Printf("issue unmarshalling cavas_state data into db.PaintEvent struct: %v\n", err)
		return
	}
	// fmt.Printf("PaintEvent: %+v\n", canvasState.canvasMap[canvasState.key][1])
	h.canvasInMemory = canvasStateMap[canvasStateKey]

}

// Load data from database into h.canvasInMemory.
// TODO_DB: Going to update this to read from the canvas_state table ... This will load a single JSON value which can be looped through to create []PaintEvent{}. Saving this will be very easy as well by adding JSONB to database.
func (h *Hub) initCanvasFromPaintEventsTable() {
	dbPaintEventsList, err := db.GetAllPaintEventsForRoom(h.roomId)
	if err != nil {
		fmt.Errorf("issue getting all paint events for room (%v): %v", h.roomId, err)
	}

	for _, dbPaintEvent := range dbPaintEventsList {
		h.canvasInMemory = append(h.canvasInMemory, &PaintEvent{
			UserId: dbPaintEvent.UserId,
			CurX:   dbPaintEvent.CurX,
			CurY:   dbPaintEvent.CurY,
			LastX:  dbPaintEvent.LastX,
			LastY:  dbPaintEvent.LastY,
			Color:  dbPaintEvent.Color,
		})
	}
}

func (h *Hub) sendActiveUserList() []byte {
	fmt.Printf("Active User List: \n")
	key := "ActiveUsers"
	userJSONList := make(map[string][]UserJSONEvents)
	// responseJSON, err := json.Marshal(h.clients)
	for _, client := range h.clients {
		userJSONList[key] = append(userJSONList[key], UserJSONEvents{
			Name:    client.user.name,
			Email:   client.user.email,
			Picture: client.user.picture,
		})
	}
	responseJSON, err := json.Marshal(userJSONList)
	if err != nil {
		fmt.Printf("get-rooms: could not create json string to return in responseText")
	}
	fmt.Printf("response: %+v\n", string(responseJSON))
	return responseJSON
}

// Start goroutine to wait numSecondsBeforeSave amout of time (in seconds) before writing the current canvas state to the canvas_state table in the database
func (h *Hub) startAutoSaveRoutine() {
	// get count of number of events in canvas memory
	prevSaveCount := len(h.canvasInMemory)

	// Start goroutine auto save ticker
	go func(prevSaveCount int) {
		for {
			select {
			case <-h.autoSave.done:
				return
			case t := <-h.autoSave.ticker.C:
				curEventCount := len(h.canvasInMemory)

				// Check if there have been any events drawn since last save, if yes, update database
				if prevSaveCount < curEventCount {
					fmt.Printf("(%v) autosave timer raached. Prev Canvas Events: %v, Cur Canvas Events: %v\n", t, prevSaveCount, curEventCount)

					h.writeCanvasStateToDB()
					prevSaveCount = curEventCount
				}
			case <-h.autoSave.event:
				curEventCount := len(h.canvasInMemory)
				if (h.autoSave.eventTriggerCount + prevSaveCount) < curEventCount {
					fmt.Printf("(%v) autosave counter reached. Prev Canvas Events: %d, Cur Canvas Events: %v\n", time.Now(), (h.autoSave.eventTriggerCount + prevSaveCount), curEventCount)

					h.writeCanvasStateToDB()
					// reset ticker to duration
					h.autoSave.ticker.Reset(h.autoSave.duration)
					prevSaveCount = curEventCount
				}
			}
		}
	}(prevSaveCount)
}

func (h *Hub) stopAutoSaveRoutine() {
	h.autoSave.done <- true
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

func getCanvasStateAsJSON(paintEventsList []*PaintEvent) ([]byte, error) {
	canvasStateMap := make(map[string][]*PaintEvent)
	canvasStateMap[canvasStateKey] = paintEventsList

	responseJSON, err := json.Marshal(canvasStateMap)
	if err != nil {
		fmt.Printf("issue Marshaling canvasStateMap\n")
		return nil, err
	}
	fmt.Printf("getCanvasStateAsJSON(): \n")
	// fmt.Printf("%s\n", responseJSON)

	return responseJSON, nil

}

func PrintHub() string {
	return "Hub"
}
