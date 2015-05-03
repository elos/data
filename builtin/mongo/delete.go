package mongo

import (
	"github.com/elos/data"
	"gopkg.in/mgo.v2"
)

func (db *MongoDB) Delete(m data.Record) error {
	s, err := db.forkSession()
	if err != nil {
		return db.err(err)
	}
	defer s.Close()

	if err = db.remove(s, m); err != nil {
		db.Printf("Error deleted record of kind %s, err: %s", m.Kind(), err)
		return err
	} else {
		db.Notify(data.NewDelete(m))
		return nil
	}
}

func (db *MongoDB) remove(s *mgo.Session, m data.Record) error {
	collection, err := db.collectionFor(s, m)
	if err != nil {
		return err
	}

	id := m.ID()

	bid, err := ParseObjectID(id.String())
	if err != nil {
		return err
	}

	if bid == emptyID {
		panic("bad id")
		//return d.ErrInvalidID
	}

	return collection.RemoveId(bid)
}
