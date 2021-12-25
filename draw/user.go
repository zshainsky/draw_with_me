package draw

import (
	"fmt"

	uuid "github.com/nu7hatch/gouuid"
)

type User struct {
	id       string   // Unique to application
	authId   string   // Map to 'sub' in JWT
	authType AuthType // Which authentication provider (i.e. Google, FB, other...)
	name     string
	email    string
	picture  string
	RoomsMap map[string]*Room // List of rooms that the user has either created or visited
}
type UserJSONEvents struct {
	Name    string `json:name`
	Email   string `json:email`
	Picture string `json:picture`
}

func NewUser(authId string, authType AuthType, name, email, picture string) *User {
	id, err := uuid.NewV4()
	if err != nil {
		fmt.Printf("problem creating unique id for client, %v", err)
		return nil
	}
	return &User{
		id:       id.String(),
		authId:   authId,
		authType: authType,
		name:     name,
		email:    email,
		picture:  picture,
		RoomsMap: make(map[string]*Room),
	}
}

func (u *User) AddRoom(room *Room) {
	fmt.Printf("user(%v) adding room: %v\n", u.email, room.Id)
	// u.RoomsList = append(u.RoomsList, room)
	// If room does not exist, add it
	if _, ok := u.RoomsMap[room.Id]; !ok {
		u.RoomsMap[room.Id] = room
	}

}

func (u *User) String() string {
	return fmt.Sprintf("User: \n\t%s\n\t%s\n\t%s\n\t%s\n\t%s\n\t%s", u.id, u.authId, u.authType, u.name, u.email, u.picture)
}

// func (u *User) isUserAdded() bool { return false }
