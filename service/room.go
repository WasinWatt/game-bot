package service

import mgo "gopkg.in/mgo.v2"

// CreateRoom creates new room
func (s *Controller) CreateRoom(session *mgo.Session, roomID string, userID string) error {
	err := s.repo.RegisterRoom(session, roomID, userID)
	if err != nil {
		return err
	}
	return nil
}
