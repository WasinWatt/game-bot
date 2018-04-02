package user

import (
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

// Register registers a user
func Register(s *mgo.Session, u *User) error {
	session := s.Copy()
	c := session.DB("undercover").C("users")
	err := c.Insert(bson.M{
		"id":     u.ID,
		"name":   u.Name,
		"roomId": u.RoomID,
	})
	return err
}

// IsExists checks if the user exists
func IsExists(s *mgo.Session, userID string) (bool, error) {
	session := s.Copy()
	c := session.DB("undercover").C("users")
	var user User
	err := c.Find(bson.M{"id": userID}).One(&user)
	if err == mgo.ErrNotFound {
		return false, nil
	}

	if err != nil {
		return false, err
	}

	return true, nil
}

// FindByRoomID finds users by room id
func FindByRoomID(s *mgo.Session, roomID string) ([]*User, error) {
	session := s.Copy()
	c := session.DB("undercover").C("users")
	var xs []*User
	err := c.Find(bson.M{"roomId": roomID}).All(&xs)
	if err != nil {
		return nil, err
	}

	return xs, nil
}
