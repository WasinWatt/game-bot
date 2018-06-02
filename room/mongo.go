package room

import (
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type Repository struct{}

func NewRepository() *Repository {
	return &Repository{}
}

// IsExists checks if room id exist
func (r *Repository) IsExists(s *mgo.Session, roomID string) (bool, error) {
	session := s.Copy()
	c := session.DB("undercover").C("rooms")
	var room Room
	err := c.Find(bson.M{"id": roomID}).One(&room)
	if err == mgo.ErrNotFound {
		return false, nil
	}

	if err != nil {
		return false, err
	}

	return true, nil
}

// Register registers a room
func (r *Repository) Register(s *mgo.Session, roomID string, userID string) error {
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
	var room Room
	err := c.Find(bson.M{"id": roomID}).One(&room)
	if err != nil {
		return false, err
	}

	if room.Owner != userID {
		return false, nil
	}

	return true, nil
}

// RemoveByID deletes room by id
func (r *Repository) RemoveByID(s *mgo.Session, roomID string) error {
	session := s.Copy()
	c := session.DB("undercover").C("rooms")
	err := c.Remove(bson.M{"id": roomID})
	if err != nil {
		return err
	}

	return nil
}
