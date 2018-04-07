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
