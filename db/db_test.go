package db

import (
	"database/sql"
	"encoding/json"
	"flag"
	"fmt"
	"testing"
	"time"
)

var userFlag bool
var roomFlag bool
var userRoomFlag bool
var paintEventFlag bool
var canvasStateFlag bool

// > go test -userFlag=true -roomFlag=true -userRoomFlag=true -paintEventFlag=true -canvasStateFlag=true
func init() {
	flag.BoolVar(&userFlag, "userFlag", true, "use to skip user db tests")
	flag.BoolVar(&roomFlag, "roomFlag", true, "use to skip room db tests")
	flag.BoolVar(&userRoomFlag, "userRoomFlag", true, "use to skip user_room db tests")
	flag.BoolVar(&paintEventFlag, "paintEventFlag", true, "use to skip paint_event db tests")
	flag.BoolVar(&canvasStateFlag, "canvasStateFlag", true, "use to skip canvas_state db tests")

}
func TestUserDB(t *testing.T) {
	t.Run("test DB connection (ping))", func(t *testing.T) {
		db := CreateConnection()
		if db == nil {
			t.Errorf("Issue connecting to DB: %+v\n", db)
		}
	})

	t.Run("get single test user", func(t *testing.T) {
		// Test user id
		userId := "62769698-ca64-472f-6da7-20becadb522a"
		got, err := GetUser(userId)

		if err != nil {
			t.Errorf("could not get user: %+v\n", err)
		}

		if got.Id != userId {
			t.Errorf("got user (%+v) is not want user (%+v)\n", got.Id, userId)
		}
	})

	t.Run("test get all users - matching row count", func(t *testing.T) {
		users, err := GetAllUsers()
		if err != nil {
			t.Error("failed on GetAllUsers()")
		}
		assertEqualTotalRowCount(t, "users", len(users))
	})

	t.Run("test insert user", func(t *testing.T) {
		// Skip this test if in Short mode (i.e: go test -short). Might want to skip to avoid altering the database
		if userFlag {
			t.Skip("skipping testing in short mode")
		}
		// Must update Id and AuthId to be unique values to make this test pass...
		testUser := UserTable{
			Id:       "62769698-ca64-472f-6da7-20becadb522c",
			AuthId:   "113062537928538714746",
			AuthType: "google",
			Name:     "test_username_2",
			Email:    "test_email@test.com",
			Picture:  "https://test.pictures.com",
		}
		got := InsertUser(testUser)

		if got != testUser.Id {
			t.Errorf("got %v, want %v", got, testUser.Id)
		}
	})

	t.Run("test update user", func(t *testing.T) {
		if userFlag {
			t.Skip("skipping testing in short mode")
		}

		testUser := UserTable{
			Id:       "62769698-ca64-472f-6da7-20becadb522c",
			AuthId:   "113062537928538714746",
			AuthType: "google",
			Name:     "test_username_3",
			Email:    "test_email@test.com",
			Picture:  "https://test.pictures.com",
		}
		got := UpdateUser(testUser)

		// should only update 1 user
		if got != 1 {
			t.Errorf("did not update user properly: %v ", got)
		}
	})

	t.Run("test delete user", func(t *testing.T) {
		if userFlag {
			t.Skip("skipping testing in short mode")
		}

		testUser := UserTable{
			Id:       "62769698-ca64-472f-6da7-20becadb522c",
			AuthId:   "113062537928538714746",
			AuthType: "google",
			Name:     "test_username_3",
			Email:    "test_email@test.com",
			Picture:  "https://test.pictures.com",
		}
		got := DeleteUser(testUser.Id)

		// should only update 1 user
		if got != 1 {
			t.Errorf("did not update user properly. Number of rows affected: %v ", got)
		}
	})
}

func TestRoomDB(t *testing.T) {
	t.Run("get single test room", func(t *testing.T) {
		// Test user id
		roomId := "62769698-ca64-472f-6da7-20becadb522b"
		got, err := GetRoom(roomId)

		if err != nil {
			t.Errorf("could not get room: %+v\n", err)
		}

		if got.Id != roomId {
			t.Errorf("got room (%+v) is not want room (%+v)\n", got.Id, roomId)
		}
	})
	t.Run("test get all rooms - matching row count", func(t *testing.T) {
		rooms, err := GetAllRooms()
		if err != nil {
			t.Error("failed on GetAllRooms()")
		}
		assertEqualTotalRowCount(t, "rooms", len(rooms))
	})
	t.Run("test insert room", func(t *testing.T) {
		// Skip this test if in Short mode (i.e: go test -short). Might want to skip to avoid altering the database
		if roomFlag {
			t.Skip("skipping testing in short mode")
		}

		// Must update Id and AuthId to be unique values to make this test pass...
		testRoom := RoomTable{
			Id:   "62769698-ca64-472f-6da7-20becadb522d",
			Name: "Paint Me A Pony!",
		}
		got := InsertRoom(testRoom)

		if got != testRoom.Id {
			t.Errorf("got %v, want %v", got, testRoom.Id)
		}
	})
	t.Run("test update room", func(t *testing.T) {
		if roomFlag {
			t.Skip("skipping testing in short mode")
		}

		testRoom := RoomTable{
			Id:   "62769698-ca64-472f-6da7-20becadb522d",
			Name: "Paint Me A Pony 2.0!",
		}
		got := UpdateRoom(testRoom)

		// should only update 1 room
		if got != 1 {
			t.Errorf("did not update room properly: %v ", got)
		}
	})
	t.Run("test delete room", func(t *testing.T) {
		if roomFlag {
			t.Skip("skipping testing in short mode")
		}

		testRoom := RoomTable{
			Id:   "62769698-ca64-472f-6da7-20becadb522d",
			Name: "Paint Me A Pony 2.0!",
		}
		got := DeleteRoom(testRoom.Id)

		// should only update 1 room
		if got != 1 {
			t.Errorf("did not update room properly. Number of rows affected: %v ", got)
		}
	})
}
func TestUserRoomDB(t *testing.T) {
	t.Run("test get all rooms for user - matching row count", func(t *testing.T) {
		userId := "62769698-ca64-472f-6da7-20becadb522a"
		got, err := GetAllRoomsForUser(userId)

		if err != nil {
			t.Errorf("issue getting all rooms for user (%v): %v", userId, err)
		}

		want := getRowCountInTableForUserOrRoom(t, "user_room", "user", userId)
		if len(got) != want {
			t.Errorf("Rooms for user (%v) did not match expected length. Got (%v), want (%v)", userId, len(got), want)
		}

	})
	t.Run("test get all users for room - matching row count", func(t *testing.T) {
		roomId := "62769698-ca64-472f-6da7-20becadb522b"
		got, err := GetAllUsersForRoom(roomId)
		if err != nil {
			t.Errorf("issue getting all user for rooms (%v): %v", roomId, err)
		}

		want := getRowCountInTableForUserOrRoom(t, "user_room", "room", roomId)
		if len(got) != want {
			t.Errorf("Users for room (%v) did not match expected length: got (%v), want (%v)", roomId, len(got), want)
		}

	})
	t.Run("test insert userroom", func(t *testing.T) {
		if userRoomFlag {
			t.Skip("skipping testing in short mode")
		}
		// Must update Id and AuthId to be unique values to make this test pass...
		testUserRoom := UserRoomTable{
			UserId: "62769698-ca64-472f-6da7-20becadb522a",
			RoomId: "62769698-ca64-472f-6da7-20becadb522f",
		}

		got := InsertUserRoom(testUserRoom)
		if testUserRoom != got {
			t.Errorf("user room could not insert user / room: got (%+v), want (%+v)", got, testUserRoom)
		}
	})

	// omitting UPDATE test

	t.Run("test delete a room for user", func(t *testing.T) {
		if userRoomFlag {
			t.Skip("skipping testing in short mode")
		}

		testUserRoom := UserRoomTable{
			UserId: "62769698-ca64-472f-6da7-20becadb522a",
			RoomId: "62769698-ca64-472f-6da7-20becadb522f",
		}

		got := DeleteUserRoom(testUserRoom.UserId, testUserRoom.RoomId)
		if got != 1 {
			t.Errorf("the desired room (%v) was not deleted for user (%v). Number of rows affected: got (%v), want (%v)", testUserRoom.UserId, testUserRoom.RoomId, got, 1)
		}
	})
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

func TestCanvasStateDB(t *testing.T) {
	t.Run("test get canvas state for room", func(t *testing.T) {
		testRoomId := "62769698-ca64-472f-6da7-20becadb522f"
		canvasState, err := GetCanvasStateForRoom(testRoomId)

		if err != nil {
			t.Errorf("issue getting canvas state for room (%v): %v", testRoomId, err)
		}

		if (canvasState == CanvasStateTable{}) {
			t.Errorf("no canvas state was returned for room (%v)", testRoomId)
			return
		}
		// epoch := time.Now().UnixNano()
		// fmt.Printf("%T\n", canvasState.CanvasJSON)
		// fmt.Printf("%d\n", epoch)

		// Test unmarshal into draw.PaintEvent
		paintEventMap := make(map[string][]*PaintEvent)

		if err := json.Unmarshal([]byte(canvasState.CanvasJSON), &paintEventMap); err != nil {
			t.Errorf("issue unmarshalling cavas_state data into db.PaintEvent struct: %v", err)
		}
		// fmt.Printf("PaintEvent: %+v\n", paintEventMap["CanvasState"][1])
	})
	t.Run("delete canvas state for room", func(t *testing.T) {
		if canvasStateFlag {
			t.Skip("skipping testing in short mode")
		}

		testRoomId := "62769698-ca64-472f-6da7-20becadb522b"

		got := DeleteCanvasStateForRoom(testRoomId)
		want := 1
		if got != int64(want) {
			t.Errorf("issue deleting canvas state for room (%v). got num rows affected (%v), want (%v)", testRoomId, got, want)
		}
	})
	t.Run("test insert canvas state for room", func(t *testing.T) {
		if canvasStateFlag {
			t.Skip("skipping testing in short mode")
		}
		testRoomId := "62769698-ca64-472f-6da7-20becadb522b"

		testCanvasState := CanvasStateTable{
			RoomId: testRoomId,
			CanvasJSON: `{
							"CanvasState": [
								{
									"CurX": 181,
									"CurY": 175,
									"LastX": 180,
									"LastY": 177,
									"Color": "#F2500F",
									"UserId": "62769698-ca64-472f-6da7-20becadb522a",
									"RoomId": "62769698-ca64-472f-6da7-20becadb522b",
									"EvtTime": 1640646814012146
								}
							]
						}`,
		}

		gotRoomId := InsertCanvasStateForRoom(testCanvasState)

		if gotRoomId != testRoomId {
			t.Errorf("got: %v, want %v", gotRoomId, testRoomId)
		}
	})

	t.Run("test update canvas state 'canvas_json' filed", func(t *testing.T) {
		if canvasStateFlag {
			t.Skip("skipping testing in short mode")
		}

		testRoomId := "62769698-ca64-472f-6da7-20becadb522b"

		testCanvasState := CanvasStateTable{
			RoomId: testRoomId,
			CanvasJSON: `{
							"CanvasState": [
								{
									"CurX": 200,
									"CurY": 200,
									"LastX": 400,
									"LastY": 400,
									"Color": "#F2500F",
									"UserId": "62769698-ca64-472f-6da7-20becadb522a",
									"RoomId": "62769698-ca64-472f-6da7-20becadb522b",
									"EvtTime": 1640646814012146
								}
							]
						}`,
		}

		gotRoomId := UpdateCanvasStateForRoom(testCanvasState)
		want := getRowCountInTableForUserOrRoom(t, "canvas_state", "room", testRoomId)
		if gotRoomId != int64(want) {
			t.Errorf("Number of rows affected does not match. got %v, want %v", gotRoomId, want)
		}
	})
}

func TestPaintEventsDB(t *testing.T) {
	t.Run("test get all paint events for room - matching row count", func(t *testing.T) {
		testRoomId := "62769698-ca64-472f-6da7-20becadb522b"
		paintEvents, err := GetAllPaintEventsForRoom(testRoomId)
		if err != nil {
			t.Errorf("issue with DB query to get all paint events: %v", err)
		}

		want := getRowCountInTableForUserOrRoom(t, "paint_event", "room", testRoomId)
		if len(paintEvents) != want {
			t.Errorf("issue getting all paint events for room (%v)", testRoomId)
		}
	})
	t.Run("test get all paint events for users - matching row count", func(t *testing.T) {
		testUserId := "62769698-ca64-472f-6da7-20becadb522a"
		paintEvents, err := GetAllPaintEventsForUser(testUserId)
		if err != nil {
			t.Errorf("issue with DB query to get all paint events: %v", err)
		}

		want := getRowCountInTableForUserOrRoom(t, "paint_event", "user", testUserId)
		if len(paintEvents) != want {
			t.Errorf("issue getting all paint events for user (%v)", testUserId)
		}
	})
	t.Run("test insert paint event", func(t *testing.T) {
		if paintEventFlag {
			t.Skip("skipping testing in short mode")
		}
		wantTime := time.Now().UTC().Unix() // INSERT using UTC
		wantUserId := "62769698-ca64-472f-6da7-20becadb522a"
		wantRoomId := "62769698-ca64-472f-6da7-20becadb522b"
		testPaintEvent := PaintEventTable{
			EvtTime: wantTime,
			UserId:  wantUserId,
			RoomId:  wantRoomId,
			LastX:   10,
			LastY:   10,
			CurX:    20,
			CurY:    20,
			Color:   "#fff000",
		}

		gotTime, gotUserId, gotRoomId := InsertPaintEvent(testPaintEvent)
		if gotTime != wantTime {
			t.Errorf("Insert paint event did not return the correct timestamp: got (%+v) want (%+v)", gotTime, wantTime)
		}
		if gotUserId != wantUserId {
			t.Errorf("Insert paint event did not return the correct userId: got (%v) want (%v)", gotUserId, wantUserId)

		}
		if gotRoomId != wantRoomId {
			t.Errorf("Insert paint event did not return the correct roomId: got (%v) want (%v)", gotRoomId, gotUserId)
		}
	})
	t.Run("test insert paint event", func(t *testing.T) {
		if paintEventFlag {
			t.Skip("skipping testing in short mode")
		}

		fmt.Println("testing all paint event insert")
		numRowsToInsert := 10
		var paintEventsList []PaintEventTable
		userId := "62769698-ca64-472f-6da7-20becadb522a"
		roomId := "62769698-ca64-472f-6da7-20becadb522b"

		for i := 0; i < numRowsToInsert; i++ {
			wantTime := time.Now().UTC().Unix() // INSERT using UTC
			paintEventsList = append(paintEventsList, PaintEventTable{
				EvtTime: wantTime,
				UserId:  userId,
				RoomId:  roomId,
				LastX:   10 + i,
				LastY:   10 + i,
				CurX:    20 + i + 10,
				CurY:    20 + i + 10,
				Color:   "#fff000",
			})
		}

		origRowCount := getTableRowCount(t, "paint_event")
		InsertAllPaintEvents(paintEventsList)

		wantNumRows := origRowCount + numRowsToInsert
		gotNumRows := getTableRowCount(t, "paint_event")

		if gotNumRows != wantNumRows {
			t.Errorf("issue on insert all paint events. Got row count (%v). Want row count (%v)", gotNumRows, wantNumRows)

		}
	})
}

func assertEqualTotalRowCount(t testing.TB, table string, numQueryResults int) {
	want := getTableRowCount(t, table)
	if numQueryResults != want {
		t.Errorf("not correct number of rows")
	}
}

// table is the table name as a string
func getTableRowCount(t testing.TB, table string) int {
	db := CreateConnection()
	defer db.Close()

	var numRows int
	// fmt.Printf("%v\n", table)

	sqlStatement := fmt.Sprintf(`SELECT count(*) as count FROM %v`, table)

	row := db.QueryRow(sqlStatement)
	err := row.Scan(&numRows)

	switch err {
	case sql.ErrNoRows:
		t.Errorf("No rows were returned!")
	case nil:
		return numRows
	default:
		t.Errorf("Unable to scan the row. %v", err)
	}
	return numRows
}
func getRowCountInTableForUserOrRoom(t testing.TB, table, idType, id string) int {
	db := CreateConnection()
	defer db.Close()

	var numRows int
	var sqlStatement string

	if idType != "room" && idType != "user" {
		t.Errorf("idType does not equal room or user. Got: %v", idType)
	}

	if idType == "room" {
		sqlStatement = fmt.Sprintf(`SELECT count(*) as count FROM %v WHERE room_id = '%v'`, table, id)
	}

	if idType == "user" {
		sqlStatement = fmt.Sprintf(`SELECT count(*) as count FROM %v WHERE user_id = '%v'`, table, id)
	}

	row := db.QueryRow(sqlStatement)
	err := row.Scan(&numRows)

	switch err {
	case sql.ErrNoRows:
		t.Errorf("No rows were returned!")
	case nil:
		return numRows
	default:
		t.Errorf("Unable to scan the row. %v", err)
	}
	return numRows
}
