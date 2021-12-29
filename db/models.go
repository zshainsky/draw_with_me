package db

type UserTable struct {
	Id       string `json:id`       // Unique to application
	AuthId   string `json:authId`   // Map to 'sub' in JWT
	AuthType string `json:authType` // Which authentication provider (i.e. Google, FB, other...)
	Name     string `json:name`
	Email    string `json:email`
	Picture  string `json:picture`
}
type UserRoomTable struct {
	UserId string `json:userId` // FK to UserTable.Id
	RoomId string `json:roomId` // FK to RoomTable.Id
}
type RoomTable struct {
	Id   string `json:id`
	Name string `json:name` // TODO: Not implemented yet in room.go
}
type CanvasStateTable struct {
	RoomId     string `json:roomId`
	CanvasJSON string `json:canvasJSON`
}
type PaintEventTable struct {
	//TODO: Add the Timestamp to the message recived from the Websocket. This is likely set in lib/components/room-canvas.js
	//TODO: Add the RoomID to the message recived from the Websocket. This is likely set in lib/components/room-canvas.js
	EvtTime int64  `json:evtTime,int64` // epoch time (ex: 1640646814.012146)
	UserId  string `json:userId`        // FK to UserTable.Id
	RoomId  string `json:roomId`        // FK to RoomTable.Id
	LastX   int    `json:lastX`
	LastY   int    `json:lastY`
	CurX    int    `json:curX`
	CurY    int    `json:curY`
	Color   string `json:color` // should be hex color including '#' (ex: #0000FF)
}
