package service

import (
	"errors"

	"github.com/WasinWatt/game-bot/user"

	mgo "gopkg.in/mgo.v2"
)

// Join joins user to the room
func (s *Controller) Join(session *mgo.Session, u *user.User) error {
	// check if user is in already in some room
	exist, err := s.repo.IsUserExists(session, u.ID)
	if err != nil {
		return err
	}
	if exist {
		return errors.New("User is already in another room. Please leave first")
	}

	// check if the room number exists
	exist, err = s.repo.IsRoomExists(session, u.RoomID)
	if err != nil {
		return err
	}
	if !exist {
		return errors.New("Room not found")
	}

	err = s.repo.RegisterUser(session, u)
	if err != nil {
		return err
	}

	return nil
}

// Leave remove user from the database
func (s *Controller) Leave(session *mgo.Session, userID string) (bool, []*user.User, error) {
	player, err := s.repo.FindUserByID(session, userID)
	if err == mgo.ErrNotFound {
		return false, nil, ErrNotFound
	}
	if err != nil {
		return false, nil, err
	}

	isOwner, err := s.repo.IsOwner(session, player.RoomID, userID)
	if err != nil {
		return false, nil, err
	}

	if isOwner {
		players, err := s.repo.FindUsersByRoomID(session, player.RoomID)
		err = s.repo.RemoveAllUsersByRoomID(session, player.RoomID)
		if err != nil {
			return true, nil, err
		}

		err = s.repo.RemoveRoomByID(session, player.RoomID)
		if err != nil {
			return true, players, err
		}

		return true, players, nil
	}

	err = s.repo.RemoveUserByID(session, userID)
	if err != nil {
		return false, nil, err
	}

	return false, nil, nil
}

// GetUser finds a user
func (s *Controller) GetUser(session *mgo.Session, userID string) (*user.User, error) {
	user, err := s.repo.FindUserByID(session, userID)
	if err == mgo.ErrNotFound {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	return user, nil
}

// GetAllUsersByRoomID finds user by room ID
func (s *Controller) GetAllUsersByRoomID(session *mgo.Session, roomID string) ([]*user.User, error) {
	xs, err := s.repo.FindUsersByRoomID(session, roomID)
	if err == mgo.ErrNotFound {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, err
	}

	return xs, nil
}

// AddUserRole add role to a user
func (s *Controller) AddUserRole(session *mgo.Session, userID string, role int) error {
	err := s.repo.SetUserRole(session, userID, role)
	if err != nil {
		return err
	}
	return nil
}
