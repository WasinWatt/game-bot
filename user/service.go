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

type Service struct {
	userdb *Repository
	roomdb *room.Repository
}

func NewService(u *Repository, r *room.Repository) *Service {
	return &Service{
		userdb: u,
		roomdb: r,
	}
}

// JoinRoom joins user to the room
func (s *Service) JoinRoom(session *mgo.Session, u *User) error {
	// check if user is in already in some room
	exist, err := s.roomdb.IsExists(session, u.ID)
	if err != nil {
		return err
	}
	if exist {
		return errors.New("User is already in another room. Please leave first")
	}

	// check if the room number exists
	exist, err = s.roomdb.IsExists(session, u.RoomID)
	if err != nil {
		return err
	}
	if !exist {
		return errors.New("Room not found")
	}

	err = s.userdb.Register(session, u)
	if err != nil {
		return err
	}

	return nil
}

// Leave remove user from the database
func (s *Service) Leave(session *mgo.Session, userID string) (bool, []*User, error) {
	player, err := s.userdb.FindByID(session, userID)
	if err == mgo.ErrNotFound {
		return false, nil, ErrNotFound
	}
	if err != nil {
		return false, nil, err
	}

	isOwner, err := s.roomdb.IsOwner(session, player.RoomID, userID)
	if err != nil {
		return false, nil, err
	}

	if isOwner {
		players, err := s.userdb.FindByRoomID(session, player.RoomID)
		err = s.userdb.RemoveAllByRoomID(session, player.RoomID)
		if err != nil {
			return true, nil, err
		}

		err = s.roomdb.RemoveByID(session, player.RoomID)
		if err != nil {
			return true, players, err
		}

		return true, players, nil
	}

	err = s.userdb.RemoveByID(session, userID)
	if err != nil {
		return false, nil, err
	}

	return false, nil, nil
}

// Get finds a user
func (s *Service) Get(session *mgo.Session, userID string) (*User, error) {
	user, err := s.userdb.FindByID(session, userID)
	if err == mgo.ErrNotFound {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	return user, nil
}

// GetAllByRoomID finds user by room ID
func (s *Service) GetAllByRoomID(session *mgo.Session, roomID string) ([]*User, error) {
	xs, err := s.userdb.FindByRoomID(session, roomID)
	if err == mgo.ErrNotFound {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, err
	}

	return xs, nil
}

// AddRole add role to a user
func (s *Service) AddRole(session *mgo.Session, userID string, role int) error {
	err := s.userdb.SetRole(session, userID, role)
	if err != nil {
		return err
	}
	return nil
}
