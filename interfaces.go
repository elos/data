package data

/*
	Model type or class
	Should correspond with the model name, generally lowercase
*/
type Kind string

// A generic interface for using data IDs
// (This is partially specific to fulfill the bson spec)
// needs work
type ID interface {
	String() string
	Hex() string
	Valid() bool
}

// Can be identied by and labeled by an ID
type Identifiable interface {
	SetID(ID)
	GetID() ID
}

/*
	To be able to be persisted something must have
	some sort of id and type
*/
type Persistable interface {
	Identifiable
	Kind() Kind
}

/*
	Relational Row
	Mongo Document

	Defines the helper function Save

	Concerned is for model update notifications
*/
type Record interface {
	Persistable

	Save(DB) error
	Concerned() []ID // for model updates
}

type ChangeKind int

const Update ChangeKind = 1
const Delete ChangeKind = 2

type Change struct {
	ChangeKind
	Record
}

func NewChange(kind ChangeKind, r Record) *Change {
	return &Change{
		ChangeKind: kind,
		Record:     r,
	}
}

/*
	Abstraction of a DataStore
	- Covers underlying connection
	- Covers id generation
	- Covers Save, Delete and simple finds
	- Covers advanced querying
	- Covers database typing for model compatability
	- Covers Registering for changeset updates
*/
type DB interface {
	// Management
	Connect(string) error

	// Persistence
	NewObjectID() ID
	CheckID(ID) error
	Save(Record) error
	Delete(Record) error
	PopulateById(Record) error
	PopulateByField(string, interface{}, Record) error

	NewQuery(Kind) Query

	Type() string

	RegisterForUpdates(Identifiable) *chan *Change
}

/* Basic, basic abstraction of underlying DBConnection */
type DBConnection interface {
	Close()
}

/* JSON attributes to values */
type AttrMap map[string]interface{}

/* data.Kinds to records */
type KindMap map[Kind]Record

/*
	Abstraction of a database query
	- Covers selection
	- Covers limiting, skipping and batching

	Executing a query returns a record iterator, the
	most elegant method of reading database lookup results
*/
type Query interface {
	Execute() (RecordIterator, error)

	Select(AttrMap) Query
	Limit(int) Query
	Skip(int) Query
	Batch(int) Query
}

/*
	Abstraction of all database query results

	Acts like an iterator - code can be written for n
	results (batching will always handle memory load)
*/
type RecordIterator interface {
	Next(Record) bool
	Close() error
}
