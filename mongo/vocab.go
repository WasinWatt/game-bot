package mongo

import (
	"github.com/WasinWatt/game-bot/vocab"
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

// RegisterVocab registers vocab
func (r *Repository) RegisterVocab(s *mgo.Session, firstWord string, secondWord string) error {
	session := s.Copy()
	c := session.DB("undercover").C("vocab")
	err := c.Insert(bson.M{"first": firstWord, "second": secondWord})
	if err != nil {
		return err
	}

	return nil
}

// CountVocab counts all data in vocab collection
func (r *Repository) CountVocab(s *mgo.Session) (int, error) {
	session := s.Copy()
	c := session.DB("undercover").C("vocab")
	len, err := c.Find(bson.M{}).Count()
	if err == mgo.ErrNotFound {
		return 0, nil
	}

	if err != nil {
		return 0, err
	}

	return len, nil
}

// FindVocabRandomly randoms one vocab
func (r *Repository) FindVocabRandomly(s *mgo.Session, skip int) (*vocab.Vocab, error) {
	session := s.Copy()
	c := session.DB("undercover").C("vocab")
	var v vocab.Vocab
	err := c.Find(bson.M{}).Skip(skip).One(&v)
	if err == mgo.ErrNotFound {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	return &v, nil
}
