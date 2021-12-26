package db

import (
	"database/sql"
	"fmt"
	"testing"
)

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
		got := len(users)
		want := getTableRowCount(t, "users")
		if got != want {
			t.Errorf("not correct number of rows")
		}
	})

	t.Run("test insert user", func(t *testing.T) {
		// Skip this test if in Short mode (i.e: go test -short). Might want to skip to avoid altering the database
		if testing.Short() {
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
		if testing.Short() {
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
		if testing.Short() {
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

	t.Run("get single test user", func(t *testing.T) {
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
}

// table is the table name as a string
func getTableRowCount(t testing.TB, table string) int {
	db := CreateConnection()
	defer db.Close()

	var numRows int
	fmt.Printf("%v\n", table)

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
