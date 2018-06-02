package vocab

import mgo "gopkg.in/mgo.v2"

type Service struct {
	vocabdb *Repository
}

func NewService(v *Repository) *Service {
	return &Service{
		vocabdb: v,
	}
}

// Add adds a new vocab
func (s *Service) Add(session *mgo.Session, firstWord string, secondWord string) error {
	err := s.vocabdb.Register(session, firstWord, secondWord)
	if err != nil {
		return err
	}

	return nil
}

// Get returns a random vocab
func (s *Service) Get(session *mgo.Session) (Vocab, error) {
	return Vocab{}, nil
}
