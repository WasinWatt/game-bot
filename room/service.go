package room

import mgo "gopkg.in/mgo.v2"

// Create creates new room
func Create(s *mgo.Session, roomID string, userID string) error {
	err := Register(s, roomID, userID)
	if err != nil {
		return err
	}
	return nil
}
