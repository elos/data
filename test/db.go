package test

import (
	"fmt"
	"github.com/elos/data"
)

type TestDB struct {
	ModelUpdates     chan *data.Change
	Saved            []data.Record
	Deleted          []data.Record
	PopulatedById    []data.Record
	PopulatedByField []data.Record

	Error bool
}

var TestDBError error = fmt.Errorf("TestDB Error")

const DBType data.DBType = "test"

type TestDBID string

func (id TestDBID) String() string {
	return string(id)
}

func (id TestDBID) Hex() string {
	return string(id)
}

func (id TestDBID) Valid() bool {
	return true
}

func NewDB() *TestDB {
	db := &TestDB{}
	db.Reset()
	return db
}

func (db *TestDB) Type() data.DBType {
	return DBType
}

func (db *TestDB) Reset() {
	db.ModelUpdates = make(chan *data.Change)
	db.Saved = make([]data.Record, 0)
	db.Deleted = make([]data.Record, 0)
	db.PopulatedById = make([]data.Record, 0)
	db.PopulatedByField = make([]data.Record, 0)
	db.Error = false
}

func (db *TestDB) Connect(addr string) error {
	if db.Error {
		return TestDBError
	}

	return nil
}

func (db *TestDB) RegisterForUpdates(a data.Identifiable) *chan *data.Change {
	return &db.ModelUpdates
}

func (db *TestDB) NewObjectID() data.ID {
	return TestDBID("")
}

func (db *TestDB) CheckID(id data.ID) error {
	if db.Error {
		return TestDBError
	}

	return nil
}

func (db *TestDB) Save(m data.Record) error {
	if db.Error {
		return TestDBError
	}

	db.Saved = append(db.Saved, m)
	return nil
}

func (db *TestDB) Delete(m data.Record) error {
	if db.Error {
		return TestDBError
	}

	db.Deleted = append(db.Deleted, m)
	return nil
}

func (db *TestDB) PopulateById(m data.Record) error {
	if db.Error {
		return TestDBError
	}

	db.PopulatedById = append(db.PopulatedById, m)
	return nil
}

func (db *TestDB) PopulateByField(field string, value interface{}, m data.Record) error {
	if db.Error {
		return TestDBError
	}

	db.PopulatedByField = append(db.PopulatedByField, m)
	return nil
}

func (db *TestDB) NewQuery(k data.Kind) data.Query {
	return nil
}
