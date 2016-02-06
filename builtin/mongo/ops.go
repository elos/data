package mongo

import (
	"github.com/elos/data"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

func (db *DB) Save(r data.Record) error {
	s, err := db.Fork()
	if err != nil {
		return err
	}
	defer s.Close()

	collection, err := db.Collection(s, r.Kind())
	if err != nil {
		return err
	}

	id := r.ID()
	bid, err := ParseObjectID(id.String())
	if err != nil {
		return data.ErrInvalidID
	}

	_, err = collection.UpsertId(bid, r)

	if err == nil {
		db.hub.Notify(data.NewUpdate(r))
	}

	return err
}

func (db *DB) Delete(r data.Record) error {
	s, err := db.Fork()
	if err != nil {
		return err
	}
	defer s.Close()

	collection, err := db.Collection(s, r.Kind())
	if err != nil {
		return err
	}

	id := r.ID()
	bid, err := ParseObjectID(id.String())
	if err != nil {
		return data.ErrInvalidID
	}

	err = collection.RemoveId(bid)

	if err == nil {
		db.hub.Notify(data.NewDelete(r))
	}

	return err
}

func (db *DB) PopulateByID(r data.Record) error {
	s, err := db.Fork()
	if err != nil {
		return err
	}
	defer s.Close()

	collection, err := db.Collection(s, r.Kind())
	if err != nil {
		return err
	}

	id := r.ID()
	bid, err := ParseObjectID(id.String())
	if err != nil {
		return data.ErrInvalidID
	}

	err = collection.FindId(bid).One(r)
	if err == mgo.ErrNotFound {
		return data.ErrNotFound
	} else {
		return err
	}
}

func (db *DB) PopulateByField(field string, value interface{}, r data.Record) error {
	s, err := db.Fork()
	if err != nil {
		return err
	}
	defer s.Close()

	collection, err := db.Collection(s, r.Kind())
	if err != nil {
		return err
	}

	err = collection.Find(bson.M{field: value}).One(r)
	if err == mgo.ErrNotFound {
		return data.ErrNotFound
	} else {
		return err
	}
}

func (db *DB) Query(k data.Kind) data.Query {
	return &Query{
		db:    db,
		kind:  k,
		match: data.AttrMap{},
	}
}
