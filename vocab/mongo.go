package vocab

import (
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type Repository struct{}

func NewRepository() *Repository {
	return &Repository{}
}

// Register registers vocab
func (r *Repository) Register(s *mgo.Session, firstWord string, secondWord string) error {
	session := s.Copy()
	c := session.DB("undercover").C("vocab")
	err := c.Insert(bson.M{"first": firstWord, "second": secondWord})
	if err != nil {
		return err
	}

	return nil
}
