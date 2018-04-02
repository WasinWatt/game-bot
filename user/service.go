package user

import (
	"errors"

	"github.com/WasinWatt/game-bot/room"
	mgo "gopkg.in/mgo.v2"
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
