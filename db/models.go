package db

import "github.com/zshainsky/draw-with-me/draw"

type UserTable struct {
	Id       string        `json:id`       // Unique to application
	AuthId   string        `json:authId`   // Map to 'sub' in JWT
	AuthType draw.AuthType `json:authType` // Which authentication provider (i.e. Google, FB, other...)
	Name     string        `json:name`
	Email    string        `json:email`
	Picture  string        `json:picture`
	// RoomsMap map[string]*draw.Room `json:roomsMap` // List of rooms that the user has either created or visited
}
type UserRoomTable struct {
	UserId string `json:userId` // FK to UserTable.Id
	RoomId string `json:roomId` // FK to RoomTable.Id
}
type RoomTable struct {
	Id   string `json:id`
	Name string `json:name` // TODO: Not implemented yet in room.go
}
type PaintEventTable struct {
	//TODO: Add the RoomID to the message recived from the Websocket. This is likely set in lib/components/room-canvas.js
	UserId string `json:userId` // FK to UserTable.Id
	RoomId string `json:roomId` // FK to RoomTable.Id
	LastX  int    `json:lastX`
	LastY  int    `json:lastY`
	CurX   int    `json:curX`
	CurY   int    `json:curY`
	Color  string `json:color` // should be hex color including '#' (ex: #0000FF)
}
