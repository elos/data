package data

import "time"

var ExampleDBType DBType = "example"
var ExampleKind Kind = "example"

type EM struct {
	Hello string
	World int
	*RecorderID
}

func NewExampleModel() *EM {
	return &EM{
		RecorderID: NewRecorderID("example"),
	}
}

// Model Constructor
func NewEM(s Store) (Model, error) {
	return NewExampleModel(), nil
}

func (em *EM) DBType() DBType {
	return ExampleDBType
}

func (em *EM) Kind() Kind {
	return ExampleKind
}

func (em *EM) ID() ID {
	return em.RecorderID
}

func (em *EM) Version() int {
	return 0
}

func (em *EM) Valid() bool {
	return true
}

func (em *EM) Concerned() []ID {
	return make([]ID, 0)
}

func (em *EM) SetID(ID) {
}

var exampleCanRead = func() bool { return true }
var exampleCanWrite = func() bool { return true }

func (em *EM) CanRead(c Client) bool {
	return exampleCanRead()
}

func (em *EM) CanWrite(c Client) bool {
	return exampleCanWrite()
}

var exampleLink = func(m Model, l Link) error { return nil }
var exampleUnlink = func(m Model, l Link) error { return nil }

func (em *EM) Link(m Model, l Link) error {
	return exampleLink(m, l)
}

func (em *EM) Unlink(m Model, l Link) error {
	return exampleUnlink(m, l)
}

func (em *EM) SetCreatedAt(t time.Time) {
}

func (em *EM) SetUpdatedAt(t time.Time) {
}

func (em *EM) UpdatedAt() time.Time {
	return time.Now()
}

func (em *EM) CreatedAt() time.Time {
	return time.Now()
}

func (em *EM) Schema() Schema {
	return NewNullSchema()
}

var ExampleLink = &Link{
	Name:    "example",
	Kind:    MulLink,
	Other:   ExampleKind,
	Inverse: "example's_inverse",
}
