package data

import (
	"time"
)

type Linker interface {
	Link(Model, Model, LinkName) error
	Unlink(Model, Model, LinkName) error
}

type Linkable interface {
	Link(Model, LinkName, Link) error
	Unlink(Model, LinkName, Link) error
}

type Validateable interface {
	Valid() bool
}

type Versioned interface {
	Version() int
}

type Schema interface {
	Linker
	Versioned
}

type Createable interface {
	CreatedAt() time.Time
	SetCreatedAt(time.Time)
}

type Updateable interface {
	UpdatedAt() time.Time
	SetUpdatedAt(time.Time)
}

type Model interface {
	Record
	Versioned
	Validateable

	Linkable
	Createable
	Updateable

	Schema() Schema
}

// === Common model patterns ===

type Nameable interface {
	Name() string
	SetName(string)
}

type Timeable interface {
	StartTime() time.Time
	SetStartTime(time.Time)
	EndTime() time.Time
	SetEndTime(time.Time)
}
