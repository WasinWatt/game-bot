package room

import (
	mgo "gopkg.in/mgo.v2"
)

type Service struct {
	roomdb *Repository
}

func NewService(r *Repository) *Service {
	return &Service{
		roomdb: r,
	}
}

// Create creates new room
func (s *Service) Create(session *mgo.Session, roomID string, userID string) error {
	err := s.roomdb.Register(session, roomID, userID)
	if err != nil {
		return err
	}
	return nil
}
