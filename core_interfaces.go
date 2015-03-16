package data

import "sync"

// ID is a generic interface for working with IDs
type ID interface {
	String() string
	Valid() bool
}

// An Identifiable can be identified by and assigned an ID
type Identifiable interface {
	ID() ID
}

type Kindable interface {
	Kind() Kind
}

type Client interface {
	Identifiable
	Kindable
}

/*
	A Persistable can be saved by a DB.

	To be able to be persisted a type must
	be identifiable and define its Kind and DBType
*/
type Persistable interface {
	Identifiable
	Kindable
	SetID(ID) error
	DBType() DBType
}

/*
	A Record represent the structured data a DB knows
	how to persist.

	It represents a relational row or mongo document

	A record must define who would be concerned
	if it were to change (underlying mechanism for subscribing
	to a DB's changes.
*/
type Record interface {
	Persistable

	Concerned() []ID // for model updates
}

type DBOps interface {
	Save(Record) error
	Delete(Record) error
	PopulateByID(Record) error
	PopulateByField(string, interface{}, Record) error
	NewQuery(Kind) RecordQuery
	RegisterForChanges(Client) *chan *Change
}

/*
	A DB is an abstraction of a database or store. A DB is anything
	that can persist and retrieve structured data, namely Records.

	A DB must handle its underlying connection,
	expose a DBType,
	implement its id generation,
	implement Save and Delete and simple Finds.

	A DB must also cover advanced selection querying,
	and provide a mechanism of subscribing to data changes.
*/
type DB interface {
	// Management
	Connect(string) error

	// Persistence
	NewID() ID
	CheckID(ID) error
	Type() DBType

	DBOps
}

/*
	A Query is an abstraction of a database query.

	It handles selection, limiting, skipping and batching.

	Executing a query returns a RecordIterator,
	an expandable and elegant method of handling n-sized
	results.
*/
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

/*
	A RecordIterator is an abstraction for handling
	Query results.

	A RecordIterator acts like an iterator - code should
	be written for n results.

	Handle memory load by batching a query's results, see: Query.Batch
*/
type RecordIterator interface {
	Next(Record) bool
	Close() error

	sync.Locker
}

type ModelIterator interface {
	Next(Model) bool
	Close() error

	sync.Locker
}
