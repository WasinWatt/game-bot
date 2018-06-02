package user

import (
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type Repository struct{}

func NewRepository() *Repository {
	return &Repository{}
}

// Register registers a user
func (r *Repository) Register(s *mgo.Session, u *User) error {
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
func (r *Repository) IsExists(s *mgo.Session, userID string) (bool, error) {
	session := s.Copy()
	c := session.DB("undercover").C("users")
	var u User
	err := c.Find(bson.M{"id": userID}).One(&u)
	if err == mgo.ErrNotFound {
		return false, nil
	}

	if err != nil {
		return false, err
	}

	return true, nil
}

// FindByRoomID finds users by room id
func (r *Repository) FindByRoomID(s *mgo.Session, roomID string) ([]*User, error) {
	session := s.Copy()
	c := session.DB("undercover").C("users")
	var xs []*User
	err := c.Find(bson.M{"roomId": roomID}).All(&xs)
	if err != nil {
		return nil, err
	}

	return xs, nil
}

// FindByID finds user by user id
func (r *Repository) FindByID(s *mgo.Session, userID string) (*User, error) {
	session := s.Copy()
	c := session.DB("undercover").C("users")
	var u User
	err := c.Find(bson.M{"id": userID}).One(&u)
	if err != nil {
		return nil, err
	}

	return &u, nil
}

// RemoveByID deletes a user by user id
func (r *Repository) RemoveByID(s *mgo.Session, userID string) error {
	session := s.Copy()
	c := session.DB("undercover").C("users")
	err := c.Remove(bson.M{"id": userID})
	if err != nil {
		return err
	}

	return nil
}

// RemoveAllByRoomID deletes all users in room id
func (r *Repository) RemoveAllByRoomID(s *mgo.Session, roomID string) error {
	session := s.Copy()
	c := session.DB("undercover").C("users")
	_, err := c.RemoveAll(bson.M{"roomId": roomID})
	if err != nil {
		return err
	}

	return nil
}

// SetRole set user role
func (r *Repository) SetRole(s *mgo.Session, userID string, role int) error {
	session := s.Copy()
	c := session.DB("undercover").C("users")
	err := c.Update(bson.M{"id": userID}, bson.M{"role": role})
	if err != nil {
		return err
	}

	return nil
}
