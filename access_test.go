package data

import (
	"testing"
	"time"
)

func TestAnonAccess(t *testing.T) {
	a := NewAnonAccess(NewNullStore())
	if a.Client() != AnonClient {
		t.Errorf("AnonAccess's client should be the AnonClient")
	}

	a2 := NewAnonAccess(NewNullStore())
	if a.Client() != a2.Client() {
		t.Errorf("All AnonAccess's should share a single AnonClient")
	}

	id := a.Client().ID()
	if id == nil {
		t.Errorf("ID should be non-nil")
	}

	kind := a.Client().Kind()
	if kind != Anonymous {
		t.Errorf("Kind should be Anonymous")
	}
}

var falsey = func() bool { return false }
var c = make(chan Record)
var m = NewNullModel()
var sendToC = func(r Record) { go func() { c <- r }() }

func TestNewAccess(t *testing.T) {
	a := NewAccess(NewNullModel(), NewRecorderStore(NewRecorderDB(), NewNullSchema()))
	if a == nil {
		t.Errorf("NewAccess should never return nil")
	}
}

func TestAccessSave(t *testing.T) {
	a := NewAccess(NewNullModel(), NewRecorderStore(NewRecorderDB(), NewNullSchema()))
	pexampleCanWrite := exampleCanWrite
	precordedSave := recordedSave
	defer func() {
		recordedSave = precordedSave
		exampleCanWrite = pexampleCanWrite
	}()

	recordedSave = sendToC

	if err := a.Save(m); err != nil {
		t.Errorf("Save should go off successfully")
	}

	select {
	case mm := <-c:
		mm, ok := mm.(*NM)
		if !ok {
			t.Errorf("Model should have been sent through c")
		}
	case <-time.After(10 * time.Millisecond):
		t.Errorf("Timed out waiting for record on recordedSave")
	}

	exampleCanWrite = falsey

	if err := a.Save(m); err != ErrAccessDenial {
		t.Errorf("Save should have returned access denial")
	}
}

func TestAccessDelete(t *testing.T) {
	a := NewAccess(NewNullModel(), NewRecorderStore(NewRecorderDB(), NewNullSchema()))
	defer func(foo func() bool, bar func(r Record)) {
		exampleCanWrite = foo
		recordedDelete = bar
	}(exampleCanWrite, recordedDelete)

	recordedDelete = sendToC

	if err := a.Delete(m); err != nil {
		t.Errorf("Delete should go off successfully")
	}

	select {
	case mm := <-c:
		mm, ok := mm.(*NM)
		if !ok {
			t.Errorf("Model should have been sent through c")
		}

		if mm.(*NM) != m {
			t.Errorf("Model sent should be the same deleted")
		}
	case <-time.After(10 * time.Millisecond):
		t.Errorf("Timed out waiting for model on delete channel")
	}

	exampleCanWrite = falsey

	if err := a.Delete(m); err != ErrAccessDenial {
		t.Errorf("Delete should now deny access")
	}

}

func TestAccessPopulateByID(t *testing.T) {
	a := NewAccess(NewNullModel(), NewRecorderStore(NewRecorderDB(), NewNullSchema()))
	if err := a.PopulateByField("foo", "bar", m); err != ErrUndefinedKind {
		t.Errorf("PopulateByField should choke on kind")
	}
	a.Register(NullKind, NewNM)
	pID := recordedPopulateByID
	pCR := exampleCanRead
	defer func() {
		recordedPopulateByID = pID
		exampleCanRead = pCR
	}()

	recordedPopulateByID = sendToC

	if err := a.PopulateByID(m); err != nil {
		t.Errorf("PopulateByID should go off successfully, got err: %s", err)
	}

	// PopulateByID sends two models, as it first checks a
	// temp, this is an optimization todo
	<-c

	select {
	case mm := <-c:
		if mm.(*NM) != m {
			t.Errorf("Model sent should be the same as populated, got %+v, wanted %+v", mm, m)
		}
	case <-time.After(10 * time.Millisecond):
		t.Errorf("Timed out waiting for model channel")
	}

	exampleCanRead = falsey

	if err := a.PopulateByID(m); err != ErrAccessDenial {
		t.Errorf("PopulateByID should now deny access, err was %s", err)
	}

	// recieve the temp model
	<-c
}

func TestAccessPopulateByField(t *testing.T) {
	a := NewAccess(NewNullModel(), NewRecorderStore(NewRecorderDB(), NewNullSchema()))

	if err := a.PopulateByField("foo", "bar", m); err != ErrUndefinedKind {
		t.Errorf("PopulateByField should choke on kind")
	}

	a.Register(NullKind, NewNM)
	defer func(foo func() bool, bar func(string, interface{}, Record)) {
		exampleCanRead = foo
		recordedPopulateByField = bar
	}(exampleCanRead, recordedPopulateByField)

	recordedPopulateByField = func(s string, i interface{}, r Record) {
		go func() { c <- r }()
	}

	if err := a.PopulateByField("foo", "bar", m); err != nil {
		t.Errorf("PopulateByField shouldn't error, got: %s", err)
	}

	// PopulateByField als sends two models
	<-c

	select {
	case mm := <-c:
		if mm.(*NM) != m {
			t.Errorf("Model sent should be the same as populated, got %+v, wanted %+v", mm, m)
		}
	case <-time.After(10 * time.Millisecond):
		t.Errorf("Timed out waiting for model channel")
	}

	exampleCanRead = falsey

	if err := a.PopulateByField("foo", "barr", m); err != ErrAccessDenial {
		t.Errorf("PopulateByField should now deny access")
	}
}

func TestAccessRegisterForUpdates(t *testing.T) {
	db := NewRecorderDB()
	a := NewAccess(NewNullModel(), NewRecorderStore(db, NewNullSchema()))

	if foo := a.RegisterForChanges(a.Client()); *foo != db.ModelUpdates {
		t.Errorf("Access register for updates failed")
	}
}

func TestAccessUnmarshal(t *testing.T) {
	a := NewAccess(NewNullModel(), NewRecorderStore(NewRecorderDB(), NewNullSchema()))
	a.Register(NullKind, NewNM)

	attrs := AttrMap{
		"string": "nick",
	}

	m, _ := a.Unmarshal(NullKind, attrs)

	if m.(*NM).String != "nick" {
		t.Errorf("Access Unmarshal failed")
	}
}
