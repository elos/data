package data

import (
	"testing"
	"time"
)

var a = NewAccess(NewExampleModel(), NewRecorderStore(NewRecorderDB(), NewNullSchema()))
var falsey = func() bool { return false }
var c = make(chan Record)
var m = NewExampleModel()
var sendToC = func(r Record) { go func() { c <- r }() }

func TestNewAccess(t *testing.T) {
	if a == nil {
		t.Errorf("NewAccess should never return nil")
	}
}

func TestAccessSave(t *testing.T) {
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
		mm, ok := mm.(*EM)
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
		mm, ok := mm.(*EM)
		if !ok {
			t.Errorf("Model should have been sent through c")
		}

		if mm.(*EM) != m {
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
	a.Store.Register(ExampleKind, NewEM)
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
		if mm.(*EM) != m {
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
	a.Store.Register(ExampleKind, NewEM)
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
		if mm.(*EM) != m {
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
