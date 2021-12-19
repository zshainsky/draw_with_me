package draw

import (
	"fmt"

	uuid "github.com/nu7hatch/gouuid"
)

type User struct {
	id        string   // Unique to application
	authId    string   // Map to 'sub' in JWT
	authType  AuthType // Which authentication provider (i.e. Google, FB, other...)
	name      string
	email     string
	picture   string
	RoomsList []*Room // List of rooms that the user has either created or visited

}

func NewUser(authId string, authType AuthType, name, email, picture string) *User {
	id, err := uuid.NewV4()
	if err != nil {
		fmt.Printf("problem creating unique id for client, %v", err)
		return nil
	}
	return &User{
		id:        id.String(),
		authId:    authId,
		authType:  authType,
		name:      name,
		email:     email,
		picture:   picture,
		RoomsList: []*Room{},
	}
}

func (u *User) AddRoom(room *Room) {
	u.RoomsList = append(u.RoomsList, room)
}
func (u *User) String() string {
	return fmt.Sprintf("User: \n\t%s\n\t%s\n\t%s\n\t%s\n\t%s\n\t%s", u.id, u.authId, u.authType, u.name, u.email, u.picture)
}

// func (u *User) isUserAdded() bool { return false }
