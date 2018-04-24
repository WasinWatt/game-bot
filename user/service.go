package user

import (
	"errors"

	"github.com/WasinWatt/game-bot/room"
	mgo "gopkg.in/mgo.v2"
)

// Errors
var (
	ErrNotFound = errors.New("user: not found")
)

// JoinRoom joins user to the room
func JoinRoom(s *mgo.Session, u *User) error {
	// check if user is in already in some room
	exist, err := IsExists(s, u.ID)
	if err != nil {
		return err
	}
	if exist {
		return errors.New("User is already in another room. Please leave first")
	}

	// check if the room number exists
	exist, err = room.IsExists(s, u.RoomID)
	if err != nil {
		return err
	}
	if !exist {
		return errors.New("Room not found")
	}

	err = Register(s, u)
	if err != nil {
		return err
	}

	return nil
}

// Leave remove user from the database
func Leave(s *mgo.Session, userID string) (bool, []*User, error) {
	player, err := FindByID(s, userID)
	if err == mgo.ErrNotFound {
		return false, nil, ErrNotFound
	}
	if err != nil {
		return false, nil, err
	}

	isOwner, err := room.IsOwner(s, player.RoomID, userID)
	if err != nil {
		return false, nil, err
	}

	if isOwner {
		players, err := FindByRoomID(s, player.RoomID)
		err = RemoveAllByRoomID(s, player.RoomID)
		if err != nil {
			return true, nil, err
		}

		err = room.RemoveByID(s, player.RoomID)
		if err != nil {
			return true, players, err
		}

		return true, players, nil
	}

	err = RemoveByID(s, userID)
	if err != nil {
		return false, nil, err
	}

	return false, nil, nil
}

// Get finds a user
func Get(s *mgo.Session, userID string) (*User, error) {
	user, err := FindByID(s, userID)
	if err == mgo.ErrNotFound {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	return user, nil
}

// GetAllByRoomID finds user by room ID
func GetAllByRoomID(s *mgo.Session, roomID string) ([]*User, error) {
	xs, err := FindByRoomID(s, roomID)
	if err == mgo.ErrNotFound {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, err
	}

	return xs, nil
}

// AddRole add role to a user
func AddRole(s *mgo.Session, userID string, role int) error {
	err := SetRole(s, userID, role)
	if err != nil {
		return err
	}
	return nil
}
