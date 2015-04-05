package data

import (
	"encoding/json"
	"sync"
)

const MemoryDBType DBType = "memory"

var defaultIDConstructor = func() ID {
	return NewNullID("memory")
}

var defaultIDParser = func(s string) (ID, error) {
	return NewNullID(s), nil
}

var defaultIDChecker = func(id ID) error {
	return nil
}

func NewMemoryDBWithType(t DBType) *MemoryDB {
	db := NewMemoryDB()
	db.dbType = t
	return db
}

func NewMemoryDB() *MemoryDB {
	return &MemoryDB{
		idConstructor: defaultIDConstructor,
		idParser:      defaultIDParser,
		idChecker:     defaultIDChecker,
		dbType:        MemoryDBType,
		records:       make(map[string]Record),
		ChangeHub:     NewChangeHub(),
	}
}

type MemoryDB struct {
	idConstructor func() ID
	idParser      func(string) (ID, error)
	idChecker     func(ID) error
	dbType        DBType

	sync.Mutex

	recordLock sync.Mutex
	records    map[string]Record

	*ChangeHub
}

func (db *MemoryDB) SetIDConstructor(f func() ID) {
	db.Lock()
	defer db.Unlock()

	db.idConstructor = f
}

func (db *MemoryDB) Connect(s string) error {
	panic("Should not connect to a memory db")
}

func (db *MemoryDB) NewID() ID {
	db.Lock()
	defer db.Unlock()

	return db.idConstructor()
}

func (db *MemoryDB) ParseID(s string) (ID, error) {
	db.Lock()
	defer db.Unlock()

	return db.idParser(s)
}

func (db *MemoryDB) CheckID(id ID) error {
	db.Lock()
	defer db.Unlock()

	return db.idChecker(id)
}

func (db *MemoryDB) Type() DBType {
	db.Lock()
	defer db.Unlock()

	return db.dbType
}

func (db *MemoryDB) Save(r Record) error {
	db.recordLock.Lock()
	defer db.recordLock.Unlock()

	if !r.ID().Valid() {
		return ErrInvalidID
	}

	db.records[r.ID().String()] = r

	db.Notify(NewUpdate(r))

	return nil
}

func (db *MemoryDB) Delete(r Record) error {
	db.recordLock.Lock()
	defer db.recordLock.Unlock()

	_, ok := db.records[r.ID().String()]
	if !ok {
		return ErrNotFound
	}

	delete(db.records, r.ID().String())

	db.Notify(NewDelete(r))

	return nil
}

func (db *MemoryDB) PopulateByID(r Record) error {
	db.recordLock.Lock()
	defer db.recordLock.Unlock()

	stored, ok := db.records[r.ID().String()]
	if !ok {
		return ErrNotFound
	}

	if bytes, err := json.Marshal(stored); err != nil {
		return err
	} else {
		return json.Unmarshal(bytes, r)
	}
}

func (db *MemoryDB) PopulateByField(name string, v interface{}, r Record) error {
	db.recordLock.Lock()
	defer db.recordLock.Unlock()

	// this is gonna be reflection
	// this sucks
	// TODO FIXME
	panic("MemoryDB PopulateByField not implemented")
	return nil
}

func (db *MemoryDB) NewQuery(k Kind) RecordQuery {
	panic("MemoryDB NewQuery not implemented")
	return nil
}
