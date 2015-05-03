package data

import "sync"

type Kind string

type ID string

func IDString(s string) ID {
	return ID(s)
}

type AttrMap map[string]interface{}

func (id ID) String() string {
	return string(id)
}

type DBType string

func (t DBType) String() string {
	return string(t)
}

type Record interface {
	ID() ID
	SetID(ID)
	Kind() Kind
	Concerned() []ID
}

type Model interface {
	Record
}

type Client interface {
	Record
}

type RecordIterator interface {
	Next(Record) bool
	Close() error

	sync.Locker
}

type RecordQuery interface {
	Execute() (RecordIterator, error)

	Select(AttrMap) RecordQuery
	Limit(int) RecordQuery
	Skip(int) RecordQuery
	Batch(int) RecordQuery
}
type ModelQuery interface {
	Execute() (ModelIterator, error)

	Select(AttrMap) ModelQuery
	Limit(int) ModelQuery
	Skip(int) ModelQuery
	Batch(int) ModelQuery
}

type ModelIterator interface {
	Next(Model) bool
	Close() error

	sync.Locker
}

type DB interface {
	// Management
	Connect(string) error

	// Persistence
	NewID() ID
	CheckID(ID) error
	ParseID(string) (ID, error)
	Type() DBType

	Save(Record) error
	Delete(Record) error
	PopulateByID(Record) error
	PopulateByField(string, interface{}, Record) error
	NewQuery(Kind) RecordQuery
	RegisterForChanges(Client) *chan *Change
}
