package data

import "sync"

// primitive type declarations
type (
	// ID is a type for indicating that a string is encodes an identifier value
	// used by the DB.
	//
	// An ID should be considered a generic form, as different database systems
	// use different mechanisms, and primitives, as IDs. Notable examples include
	// integers and UUIDs. Therefore it may be the case that an ID, which is backed
	// by a string, is actually an integer.
	//
	// Use the ID type to satisfy the Record interface.
	// Store the ID on the structure which implements the Record interface
	// as a string.
	ID string

	// Kind is a type for indicating the 'type' of a Record. For more information
	// on kinds, read about the metis data model. The Kind() method which partially
	// consitutes the Record interface allows it to be treated, in some ways, as an
	// abstract virtual class. The DB determines the 'type' by requesting the record's
	// Kind.
	//
	// Use the Kind type to satisfy the Record interface, and to
	// declare the Kinds which your program recognizes.
	Kind string

	// DBType is a type for indicating the 'type' of a DB. Similar to the Kind type,
	// the DBType allows for partial dynamic inspection of the implementation of the
	// underlying DB interface.
	//
	// Use DBType to implement the Type() method on the DB interface.
	DBType string
)

// String cases the ID as a string
func (id ID) String() string {
	return string(id)
}

// String cases the Kind as a string
func (k Kind) String() string {
	return string(k)
}

// String casts the Type as a string
func (t DBType) String() string {
	return string(t)
}

// We define types for working with persisted structures in Go.
// Specifically for Marshalling and Unmarshalling these loosely
// typed structures, with emphasis on JSON encodings.
type (
	// AttrMap is the standard type for managing the structure
	// of data records. You will note it mimics the precedent
	// set by the JSON.
	//
	// Use an AttrMap to generically deserealize JSON:
	//      bytes := []byte("{'user': { 'name': 'Nick'}})
	//		attrs := make(AttrMap)
	//		err := json.Unmarshal(bytes, &attrs)
	AttrMap map[string]interface{}

	// KindMap is the stand type for managing objects which store
	// multiple different kinds of records, and index them based on
	// their kind.
	//
	// Use a KindMap for protocols which follow
	//		{ <kind>: { ... }}
	KindMap map[Kind]interface{} //TODO should this be map[Kind]AttrMap?
)

// "The bigger the interface, the weaker the abstraction."
//
// We begin by defining the smaller constituent interfaces which we then compose into the
// larger DB interface
type (
	// An IDer manages the generation and encoding of IDs
	// for a database.
	//
	// Use the IDer interface when you require and object
	// that need only generate and parse IDs.
	IDer interface {
		NewID() ID
		ParseID(string) (ID, error)
	}

	// A Saver can persist Records from a data store.
	//
	// Use the Saver interface when you require an object
	// that need only persist records.
	Saver interface {
		Save(r Record) error
	}

	// A Deleter can remove Records from a data store.
	//
	// Use the Deleter interface when you require an object
	// that need only remove records.
	Deleter interface {
		Delete(r Record) error
	}

	// A Queryer can query Records from a data store.
	//
	// Use the Queryer interface when you require an object
	// that need only query records.
	Queryer interface {
		Query(Kind) Query
	}

	// A Populater can perform the convenience functions of doing
	// 'take first' queries on individual records.
	//
	// Use the Populater interface when you require an object that
	// need only lookup records by ID or, perhaps,  by individual
	// field names.
	Populater interface {
		PopulateByID(Record) error
		PopulateByField(string, interface{}, Record) error
	}

	DB interface {
		Type() DBType

		IDer
		Saver
		Deleter
		Populater
		Queryer

		Changes() *chan *Change
	}
)

type (
	// The most basic persistable structure
	Record interface {
		ID() ID
		SetID(ID)
		Kind() Kind
	}
)

type (
	Query interface {
		Execute() (Iterator, error)
		Skip(int) Query
		Limit(int) Query
		Batch(int) Query
		Select(AttrMap) Query
	}

	Iterator interface {
		Next(Record) bool
		Close() error
		sync.Locker
	}
)

func Equivalent(r1, r2 Record) bool {
	return r1.Kind() == r2.Kind() && r1.ID().String() == r2.ID().String()
}
