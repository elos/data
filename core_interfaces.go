package data

// ID is a generic interface for working with IDs
type ID interface {
	Valid() bool
}

// An Identifiable can be identified by and assigned an ID
type Identifiable interface {
	ID() ID
}

type Kinded interface {
	Kind() Kind
}

type Client interface {
	Identifiable
	Kinded
}

/*
	A Persistable can be saved by a DB.

	To be able to be persisted a type must
	be identifiable and define its Kind and DBType
*/
type Persistable interface {
	Identifiable
	SetID(ID)
	Kind() Kind
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
	NewQuery(Kind) Query
	RegisterForUpdates(Identifiable) *chan *Change
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

// A DBConnection serves as a basic abstraction of
// the connection underlying a DB
// FIXME: may be deprecated as the connection is
// implementation specific and the DBConnection type
// is exposed nowhere in the DB interface
type DBConnection interface {
	Close()
}

/*
	A Query is an abstraction of a database query.

	It handles selection, limiting, skipping and batching.

	Executing a query returns a RecordIterator,
	an expandable and elegant method of handling n-sized
	results.
*/
type Query interface {
	Execute() (RecordIterator, error)

	Select(AttrMap) Query
	Limit(int) Query
	Skip(int) Query
	Batch(int) Query
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
}
