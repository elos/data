package data

import (
	"encoding/json"
	"sync"
)

type MemDB struct {
	idConstructor func() ID
	idChecker     func(ID) error
	dbType        DBType

	sync.Mutex

	recordLock sync.Mutex
	records    map[ID]Record

	*ChangeHub
}

func (db *MemDB) Connect(s string) {
	panic("Should not connect to a memory db")
}

func (db *MemDB) NewID() ID {
	db.Lock()
	defer db.Unlock()

	return db.idConstructor()
}

func (db *MemDB) CheckID(id ID) error {
	db.Lock()
	defer db.Unlock()

	return db.idChecker(id)
}

func (db *MemDB) Type() DBType {
	db.Lock()
	defer db.Unlock()

	return db.dbType
}

func (db *MemDB) Save(r Record) error {
	db.recordLock.Lock()
	defer db.recordLock.Unlock()

	db.records[r.ID()] = r

	db.Notify(NewUpdate(r))

	return nil
}

func (db *MemDB) Delete(r Record) error {
	db.recordLock.Lock()
	defer db.recordLock.Unlock()

	_, ok := db.records[r.ID()]
	if !ok {
		return ErrNotFound
	}

	delete(db.records, r.ID())

	db.Notify(NewDelete(r))

	return nil
}

func (db *MemDB) PopulateByID(r Record) error {
	db.recordLock.Lock()
	defer db.recordLock.Unlock()

	stored, ok := db.records[r.ID()]
	if !ok {
		return ErrNotFound
	}

	if bytes, err := json.Marshal(stored); err != nil {
		return err
	} else {
		return json.Unmarshal(bytes, r)
	}
}

func (db *MemDB) PopulateByField(name string, v interface{}, r Record) error {
	// this sucks

	db.recordLock.Lock()
	defer db.recordLock.Unlock()

	// TODO FIXME
	panic("PopulateByField not implemented")

	return nil
}
