package mongo

import (
	"github.com/WasinWatt/game-bot/room"
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

// IsRoomExists checks if room id exist
func (r *Repository) IsRoomExists(s *mgo.Session, roomID string) (bool, error) {
	session := s.Copy()
	c := session.DB("undercover").C("rooms")
	var room room.Room
	err := c.Find(bson.M{"id": roomID}).One(&room)
	if err == mgo.ErrNotFound {
		return false, nil
	}

	if err != nil {
		return false, err
	}

	return true, nil
}

// RegisterRoom registers a room
func (r *Repository) RegisterRoom(s *mgo.Session, roomID string, userID string) error {
	session := s.Copy()
	c := session.DB("undercover").C("rooms")
	err := c.Insert(bson.M{
		"id":    roomID,
		"owner": userID,
	})
	if err != nil {
		return err
	}
	return nil
}

// IsOwner checks ownership
func (r *Repository) IsOwner(s *mgo.Session, roomID string, userID string) (bool, error) {
	session := s.Copy()
	c := session.DB("undercover").C("rooms")
	var room room.Room
	err := c.Find(bson.M{"id": roomID}).One(&room)
	if err != nil {
		return false, err
	}

	if room.Owner != userID {
		return false, nil
	}

	return true, nil
}

// RemoveRoomByID deletes room by id
func (r *Repository) RemoveRoomByID(s *mgo.Session, roomID string) error {
	session := s.Copy()
	c := session.DB("undercover").C("rooms")
	err := c.Remove(bson.M{"id": roomID})
	if err != nil {
		return err
	}

	return nil
}
