package room

import (
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

// IsExists checks if room id exist
func IsExists(s *mgo.Session, roomID string) (bool, error) {
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
func Register(s *mgo.Session, roomID string, userID string) error {
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
