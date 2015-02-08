package data

import (
	"math/rand"
	"strconv"
	"sync"
	"testing"
	"time"
)

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

func (em *EM) CanRead(c Client) bool {
	return true
}

func (em *EM) CanWrite(c Client) bool {
	return true
}

func (em *EM) Link(m Model, l Link) error {
	return nil
}

func (em *EM) Unlink(m Model, l Link) error {
	return nil
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

func TestStore(t *testing.T) {
	db := NewNullDB()
	sch := NewNullSchema()
	s := NewStore(db, sch)

	if s == nil {
		t.Errorf("NewStore should never return nil")
	}

	if len(s.RegisteredModels()) != 0 {
		t.Errorf("A new store should not have any registered models")
	}

	s.Register(ExampleKind, NewEM)

	if len(s.RegisteredModels()) != 1 || s.RegisteredModels()[0] != ExampleKind {
		t.Errorf("Register failed??")
	}

	// Concurrent registering

	c := make(chan Kind)

	var wg sync.WaitGroup
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(i int) {
			time.Sleep(time.Duration(rand.Intn(10)) * time.Millisecond)

			k := Kind(string(ExampleKind) + " " + strconv.Itoa(i))
			s.Register(k, NewEM)
			wg.Done()
			c <- k
		}(i)
	}

	wg.Wait()

	if len(s.RegisteredModels()) != 11 {
		t.Errorf("Not all the kinds were registerd, only %d were", len(s.RegisteredModels()))
	}

	kindsMap := make(map[Kind]bool)
	for _, k := range s.RegisteredModels() {
		kindsMap[k] = true
	}

	for i := 0; i < 10; i++ {
		kind := <-c
		_, ok := kindsMap[kind]
		if !ok {
			t.Errorf("%s, didn't get registered", kind)
		}
	}

	// end Concurrent registering

	// ModelFor
	m, err := s.ModelFor(ExampleKind)
	if m == nil {
		t.Errorf("ModelFor ExampleKind should give me the example model")
	}
	if err != nil {
		t.Errorf("ModelFor returned an error: %s", err)
	}
	m, ok := m.(*EM)
	if !ok {
		t.Errorf("ModelFor returned a model of the wrong type")
	}

	m, err = s.ModelFor(Kind("UnknownKind"))

	if m != nil {
		t.Errorf("ModelFor should return a nil model when unrecognized")
	}

	if err != ErrUndefinedKind {
		t.Errorf("ModelFor should return ErrUndefinedKind when model not registered")
	}
	// End ModelFor

	// Unmarshal -- this one's tough

	hello := "world"
	world := 11
	attrs := AttrMap{
		"hello": hello,
		"world": world,
	}

	m, err = s.Unmarshal(ExampleKind, attrs)

	if m == nil {
		t.Errorf("model shoul dnot be nil")
	}

	if err != nil {
		t.Errorf("Should be no error")
	}

	m, ok = m.(*EM)
	if !ok {
		t.Errorf("Unmarshal screwed up model type")
	}

	// if you've been following along, this is the amazing part
	if m.(*EM).Hello != hello || m.(*EM).World != world {
		t.Errorf("Unmarshal failed")
	}

	m, err = s.Unmarshal(Kind("Unknown Kind"), attrs)
	if err == nil {
		t.Errorf("Unmarhsall should error on an unknown kind")
	}
}
