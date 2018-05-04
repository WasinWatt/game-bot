package vocab

import mgo "gopkg.in/mgo.v2"

// Add adds a new vocab
func Add(s *mgo.Session, firstWord string, secondWord string) error {
	err := Register(s, firstWord, secondWord)
	if err != nil {
		return err
	}

	return nil
}
