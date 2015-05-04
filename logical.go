package data

import "sync"

type (
	ID     string
	Kind   string
	DBType string
)

func (id ID) String() string {
	return string(id)
}

func (k Kind) String() string {
	return string(k)
}

func (t DBType) String() string {
	return string(t)
}

type (
	// Useful for JSON based record instantiotion
	//		user.CreateWithAttrs(data.AttrMap{"name": "Nick"})
	AttrMap map[string]interface{}

	// Useful for protocols which follow
	//		{ <kind>: { ... }}
	KindMap map[Kind]interface{}
)

type (
	// The most basic persistable structure
	Record interface {
		ID() ID
		SetID(ID)
		Kind() Kind
	}

	// A slightly higher level abstraction over
	// a Record, provides for application level
	// concerns
	Model interface {
		Record
		Concerned() []ID
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

type (
	DB interface {
		Type() DBType
		NewID() ID
		ParseID(string) (ID, error)
		Save(Record) error
		Delete(Record) error
		PopulateByID(Record) error
		PopulateByField(Record) error
		NewQuery(Kind) Query
	}

	Client interface {
		Record
	}

	Access interface {
		DB
		Client() Client
	}
)
