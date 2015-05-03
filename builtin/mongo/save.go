package mongo

import (
	"github.com/elos/data"
	"gopkg.in/mgo.v2"
)

func (db *MongoDB) Save(m data.Record) error {
	s, err := db.forkSession()
	if err != nil {
		return db.err(err)
	}
	defer s.Close()

	if err = db.save(s, m); err != nil {
		db.Printf("Error saving record of kind %s, err: %s", m.Kind(), err.Error())
		return err
	} else {
		db.Notify(data.NewUpdate(m))
		return nil
	}
}

func (db *MongoDB) save(s *mgo.Session, r data.Record) error {
	collection, err := db.collectionFor(s, r)
	if err != nil {
		return err
	}

	id := r.ID()

	bid, err := ParseObjectID(id.String())
	if err != nil {
		return err
	}

	_ /*changeInfo*/, err = collection.UpsertId(bid, r)

	return err
}
