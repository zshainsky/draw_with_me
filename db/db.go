// Code Reference: https://codesource.io/build-a-crud-application-in-golang-with-postgresql/
package db

import (
	"database/sql"
	"fmt"
	"log"
	"os"

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

/* END - CRUD Functions for UserRoomTable */

/* START - CRUD Functions for UserRoomTable */
func GetAllPaintEventsForRoom(roomId string) ([]PaintEventTable, error) {
	return []PaintEventTable{}, nil
}
func GetAllPaintEventsForUser(userId string) ([]PaintEventTable, error) {
	return []PaintEventTable{}, nil
}

// string retuned will be the id of the user from the database
func InsertPaintEvent(paintEvent PaintEventTable) string {
	return ""
}
func InsertAllPaintEvents(paintEvent []PaintEventTable) string {
	return ""
}

// omitting UPDATE

// return number of rows affected
func DeletePaintEventsForUser(userId string) int64 {
	return 0
}
func DeletePaintEventsForRoom(roomId string) int64 {
	return 0
}
