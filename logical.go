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
		ParseID(encodedID string) (ID, error)
	}

	// A Saver can persist Records from a data store.
	//
	// Use the Saver interface when you require an object
	// that need only persist records.
	Saver interface {
		// Save persists a record. Save should implement the functionality of
		// a so-called 'upsert'. If the record exists, it should be updated
		// to the new structure, otherwise it should be created, and then stored.
		//
		// The meaning of 'persist' can vary over different implementations.
		// For an in memory data store, it may mean throwing the record in a
		// a map, although for other datastore it may mean encoding the Record
		// and then sending it over the wire to the database management system.
		//
		// Save may return the following errors:
		//  * ErrNoConnection
		//		- The Saver has lost connection
		//  * ErrInvalidID
		//		- The Record's ID has an invalid encoding
		//	* ErrAccessDenial
		//		- The client does not have permission to modify/create the record
		Save(r Record) error
	}

	// A Deleter can remove Records from a data store.
	//
	// Use the Deleter interface when you require an object
	// that need only remove records.
	Deleter interface {
		// Save persists a record. Save should implement the functionality of
		// of completely removing a record from the data store.
		//
		// The meaning of delete can vary from application to application, and
		// from database to database. Here we adopt the traditional definition
		// of completely erasing the structure. Were you to want to implement
		// a custom form of deletion, such as setting some attribute on the
		// structure, you could do so by implementing a structure which also
		// implemented the Deleter interface, but took special measures. However
		// the data.Deleter interface always completely erases the record. So you
		// should define your own interface
		//		package custom
		//
		//
		//		// SoftDeleter 'deletes' a record by setting it's DeletedAt
		//		// attribute.
		//		type SoftDeleter interface{
		//			data.Deleter
		//		}
		//
		// Delete may return the following errors:
		//  * ErrNoConnection
		//		- The Saver has lost connection
		//  * ErrInvalidID
		//		- The Record's ID has an invalid encoding
		//	* ErrAccessDenial
		//		- The client does not have permission to delete the record
		//
		// Note: ErrNotFound is not a valid error. If the record does not exist
		// then the Deleter should ignore the request.
		Delete(r Record) error // TODO: move to custom error type Error
	}

	// A Queryer can produce Queries which can query over
	// Records from a data store.
	//
	// Use the Queryer interface when you require an object
	// that need only query records.
	Queryer interface {
		// Query creates a new Query over the domain of the
		// given kind.
		//
		// A Queryer could also be considered a QueryProducer,
		// as it does not actually execute the query, that is
		// handled by a separate interface, namely a Query.
		Query(k Kind) Query
	}

	// A Populater can perform the convenience functions of doing
	// 'take first' queries on individual records.
	//
	// Use the Populater interface when you require an object that
	// need only lookup records by ID or, perhaps, by individual
	// field names.
	Populater interface {
		// PopulateByID populates the structure of a record by using
		// the Record's Kind() and ID().
		//
		// PopulateByID may return the following errors
		//	* ErrNotFound
		//		* The record with the given kind and id does not exist
		//  * ErrNoConnection
		//		- The Populater has lost connection
		//  * ErrInvalidID
		//		- The Record's ID has an invalid encoding
		//	* ErrAccessDenial
		//		- The client does not have permission to load the record
		//			(see access note below)
		PopulateByID(Record) error

		// PopulateByField populates the structure of a records by using
		// the record's Kind(), the field and the value. It is possible
		// that multiple records contain the (field, value) pair, in this
		// scenario the one that gets populated into the Record structure
		// is undefined. Furthermore, it is _not_ guaranteed that
		// PopulateByField should populate the same record twice.
		//
		// PopulateByField is the generic form of PopulateByID. It should
		// only be used when the field is some sort of identification, and is
		// guaranteed to be unique across the domain given by the record's kind.
		//
		// PopulateByField may return the following errors
		//	* ErrNotFound
		//		* The record with the given kind, field, and value does not exist
		//  * ErrNoConnection
		//		- The Populater has lost connection
		//	* ErrAccessDenial
		//		- The client does not have permission to load the record
		//			(see access note below)
		PopulateByField(field string, value interface{}, r Record) error

		// ErrAccessDenial Note:
		//	Returning ErrAccessDenial to an end user is poor access
		//	control, as	it leaks information, namely, that the Record
		//	does indeed exist. Therefore a client using a Populater
		// 	should, perhaps, change a ErrAccessDenial to ErrNotFound in
		//	order to avoid leaking information. However, a Populater's
		//	PopulateByID and PopulateByField methods are defined to expose
		//	the information associated 	with ErrAccessDenial)
	}

	// A DB is the composition of the individual interfaces which
	// a data store conventionally satisfies. DB is a relatively
	// high-level interface, and contains many methods, which does
	// weaken the abstraction. Still, DB maintains a layer of abstraction
	// between client programs and the management of persistent state
	//
	// We define persistent state generally to be state which survives
	// across processes and nodes, however a DB could be backed by
	// an in-memory store. Such interchangeability is indeed the benefit
	// of the DB interface.
	//
	// Use a DB to interface with persistent application state.
	DB interface {
		// Type() returns the DBType of this database.
		//
		// Use Type() conservatively, and only for dynamic inspection of
		// the database. A common use case is in printing error messages
		// and inspecting a bit more the type of DB you have
		//
		// TODO: DEPRECATE
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
