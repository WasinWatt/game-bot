package mongo

import (
	"github.com/WasinWatt/game-bot/user"
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

// RegisterUser registers a user
func (r *Repository) RegisterUser(s *mgo.Session, u *user.User) error {
	session := s.Copy()
	c := session.DB("undercover").C("users")
	err := c.Insert(bson.M{
		"id":     u.ID,
		"name":   u.Name,
		"roomId": u.RoomID,
	})
	return err
}

// IsUserExists checks if the user exists
func (r *Repository) IsUserExists(s *mgo.Session, userID string) (bool, error) {
	session := s.Copy()
	c := session.DB("undercover").C("users")
	var u user.User
	err := c.Find(bson.M{"id": userID}).One(&u)
	if err == mgo.ErrNotFound {
		return false, nil
	}

	if err != nil {
		return false, err
	}

	return true, nil
}

// FindUsersByRoomID finds users by room id
func (r *Repository) FindUsersByRoomID(s *mgo.Session, roomID string) ([]*user.User, error) {
	session := s.Copy()
	c := session.DB("undercover").C("users")
	var xs []*user.User
	err := c.Find(bson.M{"roomId": roomID}).All(&xs)
	if err != nil {
		return nil, err
	}

	return xs, nil
}

// FindUserByID finds user by user id
func (r *Repository) FindUserByID(s *mgo.Session, userID string) (*user.User, error) {
	session := s.Copy()
	c := session.DB("undercover").C("users")
	var u user.User
	err := c.Find(bson.M{"id": userID}).One(&u)
	if err != nil {
		return nil, err
	}

	return &u, nil
}

// RemoveUserByID deletes a user by user id
func (r *Repository) RemoveUserByID(s *mgo.Session, userID string) error {
	session := s.Copy()
	c := session.DB("undercover").C("users")
	err := c.Remove(bson.M{"id": userID})
	if err != nil {
		return err
	}

	return nil
}

// RemoveAllUsersByRoomID deletes all users in room id
func (r *Repository) RemoveAllUsersByRoomID(s *mgo.Session, roomID string) error {
	session := s.Copy()
	c := session.DB("undercover").C("users")
	_, err := c.RemoveAll(bson.M{"roomId": roomID})
	if err != nil {
		return err
	}

	return nil
}

// SetUserRole set user role
func (r *Repository) SetUserRole(s *mgo.Session, userID string, role int) error {
	session := s.Copy()
	c := session.DB("undercover").C("users")
	err := c.Update(bson.M{"id": userID}, bson.M{"role": role})
	if err != nil {
		return err
	}

	return nil
}
