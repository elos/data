package data

// DBType is defined by each implementation of a DB,
// so that the DB can be semi-identifiable and implement
// compatibility checking against a Record's declared DBType
type DBType string

// Kind is a record's table name of collection name.
// It should correspond to the model's name, generally lowercase.
type Kind string

// A KindMap represents a mapping from
// Record Kind to Record object, useful
// for protocols which follow:
//	{ <kind>: { ... info ... } }
type KindMap map[Kind]Record

// An AttrMap is the type used to
// populate a Record's fields.
type AttrMap map[string]interface{}

// A ChangeKind represents the nature of a Change.
type ChangeKind int

const (
	// Update is the ChangeKind triggered on a save.
	// As such  it covers boths creation and modification.
	Update ChangeKind = 1

	// Delete is the ChangeKind triggered on a delete.
	Delete ChangeKind = 2
)

// A Change represents a modification to the data state
// a DB represents. Any succesful modification to the underlying
// should trigger a Change to be sent of a channel.
// Implementations of a the DB interface should implement all
// defined ChangeKinds
type Change struct {
	ChangeKind
	Record
}

// NewChange is a simple constructor for a Change object
func NewChange(kind ChangeKind, r Record) *Change {
	return &Change{
		ChangeKind: kind,
		Record:     r,
	}
}
