// Code Reference: https://codesource.io/build-a-crud-application-in-golang-with-postgresql/
package db

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"strings"

	_ "github.com/lib/pq" // driver import
)

// response format ... may not need
type Response struct {
	ID      int64  `json:"id,omitempty"`
	Message string `json:"message,omitempty"`
}

type ConnectionParameters struct {
	user         string
	password     string
	host         string
	port         string
	databaseName string
	sslMode      string
}

func CreateConnection() *sql.DB {
	// Check env file for DATABASE_URL. This is the Heroku env name. If not found use the default.
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		dbParam := ConnectionParameters{
			user:         "postgres",
			password:     "12345",
			host:         "localhost",
			port:         "5432",
			databaseName: "draw_with_me",
			sslMode:      "disable",
		}

		dbURL = fmt.Sprintf("postgresql://%v:%v@%v:%v/%v?sslmode=%v", dbParam.user, dbParam.password, dbParam.host, dbParam.port, dbParam.databaseName, dbParam.sslMode) // Default database url if not specified
	}
	// Open the connection
	db, err := sql.Open("postgres", dbURL)

	if err != nil {
		panic(err)
	}

	// check the connection
	err = db.Ping()

	if err != nil {
		panic(err)
	}

	// fmt.Println("Successfully connected!")
	// return the connection
	return db
}

/* START - CRUD Functions for UserTable */
func GetUser(id string) (UserTable, error) {
	// connect to db and close on completion
	db := CreateConnection()
	defer db.Close()

	// create UserTable variable to load data into
	var user UserTable

	sqlStatement := `SELECT * FROM users WHERE id=$1`

	row := db.QueryRow(sqlStatement, id)
	err := row.Scan(&user.Id, &user.AuthId, &user.AuthType, &user.Name, &user.Email, &user.Picture)

	switch err {
	case sql.ErrNoRows:
		fmt.Println("No rows were returned!")
		return user, nil
	case nil:
		return user, nil
	default:
		log.Fatalf("Unable to scan the row. %v", err)
	}

	// return empty user on error
	return user, err
}
func GetAllUsers() ([]UserTable, error) {
	db := CreateConnection()
	defer db.Close()

	var users []UserTable

	sqlStatement := `SELECT * FROM users`
	rows, err := db.Query(sqlStatement)
	if err != nil {
		log.Fatalf("unable to execute query for all users. \n\tsql: (%v)\n\terror: %v", sqlStatement, err)
	}
	defer rows.Close()

	for rows.Next() {
		var user UserTable

		err := rows.Scan(&user.Id, &user.AuthId, &user.AuthType, &user.Name, &user.Email, &user.Picture)
		if err != nil {
			log.Fatalf("unable to scan single row for all users query: %v", err)
		}
		users = append(users, user)
	}

	return users, err
}

// string retuned will be the id of the user from the database
func InsertUser(user UserTable) string {
	db := CreateConnection()
	defer db.Close()

	var id string
	sqlStatement := "INSERT INTO users (id, authId, authType, name, email, picture) VALUES ($1, $2, $3, $4, $5, $6) RETURNING id;"

	row := db.QueryRow(sqlStatement, user.Id, user.AuthId, user.AuthType, user.Name, user.Email, user.Picture)
	err := row.Scan(&id)
	if err != nil {
		log.Fatalf("could not insert user: %v", err)
	}

	return id
}

// return number of rows affected
func UpdateUser(user UserTable) int64 {
	db := CreateConnection()
	defer db.Close()

	sqlStatement := "UPDATE users SET AuthId=$2, AuthType=$3, Name=$4, Email=$5, Picture=$6 WHERE id=$1;"

	res, err := db.Exec(sqlStatement, user.Id, user.AuthId, user.AuthType, user.Name, user.Email, user.Picture)
	if err != nil {
		log.Fatalf("unable to update user (%v): %v", user.Id, err)
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		log.Fatalf("error while checking the affected rows. %v", err)
	}

	return rowsAffected
}

// return number of rows affected
func DeleteUser(id string) int64 {
	db := CreateConnection()
	defer db.Close()

	sqlStatement := "DELETE FROM users WHERE id=$1;"

	res, err := db.Exec(sqlStatement, id)
	if err != nil {
		log.Fatalf("unable to update user (%v): %v", id, err)
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		log.Fatalf("error while checking the affected rows. %v", err)
	}

	return rowsAffected
}

/* END - CRUD Functions for UserTable */

/* START - CRUD Functions for RoomTable */
func GetRoom(id string) (RoomTable, error) {
	// connect to db and close on completion
	db := CreateConnection()
	defer db.Close()

	// create UserTable variable to load data into
	var room RoomTable

	sqlStatement := `SELECT * FROM rooms WHERE id=$1`

	row := db.QueryRow(sqlStatement, id)
	err := row.Scan(&room.Id, &room.Name)

	switch err {
	case sql.ErrNoRows:
		fmt.Println("No rows were returned!")
		return room, nil
	case nil:
		return room, nil
	default:
		log.Fatalf("Unable to scan the row. %v", err)
	}

	// return empty user on error
	return room, err
}
func GetAllRooms() ([]RoomTable, error) {
	db := CreateConnection()
	defer db.Close()

	var rooms []RoomTable

	sqlStatement := `SELECT * FROM rooms`
	rows, err := db.Query(sqlStatement)

	if err != nil {
		log.Fatalf("unable to execute query for all users. \n\tsql: (%v)\n\terror: %v", sqlStatement, err)
	}
	defer rows.Close()

	for rows.Next() {
		var room RoomTable

		err := rows.Scan(&room.Id, &room.Name)
		if err != nil {
			log.Fatalf("unable to scan single row for all users query: %v", err)
		}
		rooms = append(rooms, room)
	}

	return rooms, err
}

// string retuned will be the id of the user from the database
func InsertRoom(room RoomTable) string {
	db := CreateConnection()
	defer db.Close()

	var id string
	sqlStatement := "INSERT INTO rooms (id, name) VALUES ($1, $2) RETURNING id;"

	row := db.QueryRow(sqlStatement, room.Id, room.Name)
	err := row.Scan(&id)
	if err != nil {
		log.Fatalf("could not insert room: %v", err)
	}

	return id
}

// return number of rows affected
func UpdateRoom(room RoomTable) int64 {
	db := CreateConnection()
	defer db.Close()

	sqlStatement := "UPDATE rooms SET Name=$2 WHERE id=$1;"

	res, err := db.Exec(sqlStatement, room.Id, room.Name)
	if err != nil {
		log.Fatalf("unable to update user (%v): %v", room.Id, err)
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		log.Fatalf("error while checking the affected rows. %v", err)
	}

	return rowsAffected
}

// return number of rows affected
func DeleteRoom(id string) int64 {
	db := CreateConnection()
	defer db.Close()

	sqlStatement := "DELETE FROM rooms WHERE id=$1;"

	res, err := db.Exec(sqlStatement, id)
	if err != nil {
		log.Fatalf("unable to update room (%v): %v", id, err)
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		log.Fatalf("error while checking the affected rows. %v", err)
	}

	return rowsAffected
}

/* END - CRUD Functions for RoomTable */

/* START - CRUD Functions for UserRoomTable */
func GetAllRoomsForUser(userId string) ([]RoomTable, error) {
	db := CreateConnection()
	defer db.Close()

	var rooms []RoomTable

	sqlStatement := `SELECT * FROM rooms where id IN (
						SELECT room_id FROM user_room WHERE user_id = $1
					);`
	rows, err := db.Query(sqlStatement, userId)

	if err != nil {
		log.Fatalf("unable to execute query for all users. \n\tsql: (%v)\n\terror: %v", sqlStatement, err)
	}
	defer rows.Close()

	for rows.Next() {
		var room RoomTable

		err := rows.Scan(&room.Id, &room.Name)
		if err != nil {
			log.Fatalf("unable to scan single row for all users query: %v", err)
		}
		rooms = append(rooms, room)
	}
	return rooms, nil
}
func GetAllUsersForRoom(roomId string) ([]UserTable, error) {
	db := CreateConnection()
	defer db.Close()

	var users []UserTable

	sqlStatement := `SELECT * FROM users where id IN (
						SELECT user_id FROM user_room WHERE room_id = $1
					);`
	rows, err := db.Query(sqlStatement, roomId)

	if err != nil {
		log.Fatalf("unable to execute query for all users. \n\tsql: (%v)\n\terror: %v", sqlStatement, err)
	}
	defer rows.Close()

	for rows.Next() {
		var user UserTable

		err := rows.Scan(&user.Id, &user.AuthId, &user.AuthType, &user.Name, &user.Email, &user.Picture)
		if err != nil {
			log.Fatalf("unable to scan single row for all users query: %v", err)
		}
		users = append(users, user)
	}
	return users, nil
}

// string retuned will be the id of the user from the database
func InsertUserRoom(userRoom UserRoomTable) UserRoomTable {
	db := CreateConnection()
	defer db.Close()

	var userId, roomId string
	sqlStatement := "INSERT INTO user_room (user_id, room_id) VALUES ($1, $2) RETURNING user_id, room_id;"

	row := db.QueryRow(sqlStatement, userRoom.UserId, userRoom.RoomId)
	err := row.Scan(&userId, &roomId)

	if err != nil {
		log.Fatalf("could not insert user: %v", err)
	}

	return UserRoomTable{UserId: userId, RoomId: roomId}
}

// omitting UPDATE

// return number of rows affected
func DeleteUserRoom(userId, roomId string) int64 {
	db := CreateConnection()
	defer db.Close()

	sqlStatement := "DELETE FROM user_room WHERE user_id=$1 and room_id=$2;"

	res, err := db.Exec(sqlStatement, userId, roomId)
	if err != nil {
		log.Fatalf("unable to delete room (%v) for user %v: %v", roomId, userId, err)
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		log.Fatalf("error while checking the affected rows. %v", err)
	}

	return rowsAffected
}

/* START - CRUD Functions for CanvasStateTable */
func GetCanvasStateForRoom(roomId string) (CanvasStateTable, error) {
	db := CreateConnection()
	defer db.Close()

	// create UserTable variable to load data into
	var canvasState CanvasStateTable

	sqlStatement := `SELECT * FROM canvas_state WHERE room_id=$1`

	row := db.QueryRow(sqlStatement, roomId)
	err := row.Scan(&canvasState.RoomId, &canvasState.CanvasJSON)

	switch err {
	case sql.ErrNoRows:
		fmt.Println("No rows were returned!")
		return canvasState, nil
	case nil:
		return canvasState, nil
	default:
		log.Fatalf("Unable to scan the row. %v", err)
	}

	// return empty user on error
	return canvasState, err
}

func InsertCanvasStateForRoom(canvasState CanvasStateTable) string {
	db := CreateConnection()
	defer db.Close()

	var id string
	sqlStatement := "INSERT INTO canvas_state (room_id, canvas_json) VALUES ($1, $2) RETURNING room_id;"

	row := db.QueryRow(sqlStatement, canvasState.RoomId, canvasState.CanvasJSON)
	err := row.Scan(&id)
	if err != nil {
		log.Fatalf("could not insert room: %v", err)
	}

	return id
}
func UpdateCanvasStateForRoom(canvasState CanvasStateTable) int64 {
	db := CreateConnection()
	defer db.Close()

	sqlStatement := "UPDATE canvas_state SET canvas_json=$2 WHERE room_id=$1;"

	res, err := db.Exec(sqlStatement, canvasState.RoomId, canvasState.CanvasJSON)
	if err != nil {
		log.Fatalf("unable to update canvas state for roomId (%v): %v", canvasState.RoomId, err)
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		log.Fatalf("error while checking the affected rows. %v", err)
	}

	return rowsAffected
}
func DeleteCanvasStateForRoom(roomId string) int64 {
	db := CreateConnection()
	defer db.Close()

	sqlStatement := "DELETE FROM canvas_state WHERE room_id=$1;"

	res, err := db.Exec(sqlStatement, roomId)
	if err != nil {
		log.Fatalf("unable to delete room (%v): %v", roomId, err)
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		log.Fatalf("error while checking the affected rows. %v", err)
	}

	return rowsAffected
}

/* END - CRUD Functions for CanvasStateTable */

/* START - CRUD Functions for PaintEventsTable */
func GetAllPaintEventsForRoom(roomId string) ([]PaintEventTable, error) {
	db := CreateConnection()
	defer db.Close()

	var paintEventsList []PaintEventTable

	sqlStatement := `SELECT * FROM paint_event WHERE room_id=$1`
	rows, err := db.Query(sqlStatement, roomId)
	if err != nil {
		log.Fatalf("unable to execute query for all users. \n\tsql: (%v)\n\terror: %v", sqlStatement, err)
	}
	defer rows.Close()

	for rows.Next() {
		var paintEvent PaintEventTable

		err := rows.Scan(&paintEvent.EvtTime, &paintEvent.UserId, &paintEvent.RoomId, &paintEvent.LastX, &paintEvent.LastY, &paintEvent.CurX, &paintEvent.CurY, &paintEvent.Color)
		if err != nil {
			log.Fatalf("unable to scan single row for all users query: %v", err)
		}
		// fmt.Printf("timestamp: (%v) %v\n", paintEvent.EvtTime.Hour(), paintEvent.EvtTime)
		paintEventsList = append(paintEventsList, paintEvent)
	}

	return paintEventsList, err
}
func GetAllPaintEventsForUser(userId string) ([]PaintEventTable, error) {
	db := CreateConnection()
	defer db.Close()

	var paintEventsList []PaintEventTable

	sqlStatement := `SELECT * FROM paint_event WHERE user_id=$1`
	rows, err := db.Query(sqlStatement, userId)
	if err != nil {
		log.Fatalf("unable to execute query for all users. \n\tsql: (%v)\n\terror: %v", sqlStatement, err)
	}
	defer rows.Close()

	for rows.Next() {
		var paintEvent PaintEventTable

		err := rows.Scan(&paintEvent.EvtTime, &paintEvent.UserId, &paintEvent.RoomId, &paintEvent.LastX, &paintEvent.LastY, &paintEvent.CurX, &paintEvent.CurY, &paintEvent.Color)
		if err != nil {
			log.Fatalf("unable to scan single row for all users query: %v", err)
		}
		paintEventsList = append(paintEventsList, paintEvent)
	}

	return paintEventsList, err
}

// string retuned will be the id of the user from the database
func InsertPaintEvent(paintEvent PaintEventTable) (time int64, userId, roomId string) {
	db := CreateConnection()
	defer db.Close()

	sqlStatement := "INSERT INTO paint_event (evt_time, user_id, room_id, lastX, lastY, curX, curY, color) VALUES ($1, $2, $3, $4, $5, $6, $7, $8) RETURNING evt_time, user_id, room_id;"

	row := db.QueryRow(sqlStatement, &paintEvent.EvtTime, &paintEvent.UserId, &paintEvent.RoomId, &paintEvent.LastX, &paintEvent.LastY, &paintEvent.CurX, &paintEvent.CurY, &paintEvent.Color)
	err := row.Scan(&time, &userId, &roomId)

	if err != nil {
		log.Fatalf("could not insert user: %v", err)
	}

	return time, userId, roomId
}

// Reference: https://github.com/TrinhTrungDung/go-bulk-create/blob/master/main_test.go
func InsertAllPaintEvents(paintEventsList []PaintEventTable) {
	db := CreateConnection()
	defer db.Close()

	tx, err := db.Begin()
	if err != nil {
		log.Fatalf("could not begin db in InsertAllPaintEvents(): %v", err)
	}

	valueStrings := []string{}
	valueArgs := []interface{}{}
	for i, paintEvent := range paintEventsList {
		position := i * 8
		valueStrings = append(valueStrings, fmt.Sprintf("($%d, $%d, $%d, $%d, $%d, $%d, $%d, $%d)", position+1, position+2, position+3, position+4, position+5, position+6, position+7, position+8))
		valueArgs = append(valueArgs, paintEvent.EvtTime)
		valueArgs = append(valueArgs, paintEvent.UserId)
		valueArgs = append(valueArgs, paintEvent.RoomId)
		valueArgs = append(valueArgs, paintEvent.CurX)
		valueArgs = append(valueArgs, paintEvent.CurY)
		valueArgs = append(valueArgs, paintEvent.LastX)
		valueArgs = append(valueArgs, paintEvent.LastY)
		valueArgs = append(valueArgs, paintEvent.Color)
	}

	sqlStatement := fmt.Sprintf("INSERT INTO paint_event (evt_time, user_id, room_id, lastX, lastY, curX, curY, color) VALUES %s;", strings.Join(valueStrings, ","))
	_, err = tx.Exec(sqlStatement, valueArgs...)
	if err != nil {
		fmt.Println("rolling back")
		tx.Rollback()
		fmt.Println(err)
	}
	err = tx.Commit()
	if err != nil {
		fmt.Println("commit error")
		fmt.Println(err)
	}
}

// omitting UPDATE

// omitting DELETE
