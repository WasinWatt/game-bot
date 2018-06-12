package service

import (
	"math/rand"

	"github.com/WasinWatt/game-bot/vocab"
	mgo "gopkg.in/mgo.v2"
)

// AddVocab adds a new vocab
func (s *Controller) AddVocab(session *mgo.Session, firstWord string, secondWord string) error {
	err := s.repo.RegisterVocab(session, firstWord, secondWord)
	if err != nil {
		return err
	}

	return nil
}

// GetVocab returns a random vocab
func (s *Controller) GetVocab(session *mgo.Session) (*vocab.Vocab, error) {
	len, err := s.repo.CountVocab(session)
	if err != nil {
		return nil, err
	}

	skip := rand.Intn(len)
	v, err := s.repo.FindVocabRandomly(session, skip)
	if err != nil {
		return nil, err
	}

	return v, nil
}
