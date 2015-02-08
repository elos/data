package data

import (
	"errors"
	"testing"
)

func TestRecorderID(t *testing.T) {
	var id = NewRecorderID("foobar")

	if id == nil {
		t.Errorf("NewRecorderID should never return nil")
	}

	if id.Valid() != true {
		t.Errorf("A new RecorderID should start as valid")
	}

	id.SetValidity(false)
	if id.Valid() != false {
		t.Errorf("SetValidity Failed, wanted false got true")
	}

	id.SetValidity(true)
	if id.Valid() != true {
		t.Errorf("SetValidity Failed, wanted true got false")
	}
}

func TestRecorderDB(t *testing.T) {
	db := NewRecorderDB()
	if !testStartRecorderDBCleanState(t, db) {
		t.Errorf("TestRecorder had bad intitial state")
	}

	a := "http://localhost:8080"
	err := db.Connect(a)
	if err != nil {
		t.Errorf("Connect() should not return an error")
	}
	if db.Connected != a {
		t.Errorf("Connect() should set db.Connected field")
	}

	id := db.NewID()
	id, ok := id.(*RecorderID)
	if !ok {
		t.Errorf("NewID should return a *RecorderID")
	}
	if id == nil {
		t.Errorf("NewID should not return nil")
	}

	err = db.CheckID(id)
	if err != nil {
		t.Errorf("CheckID should have passed, id was valid")
	}

	id.(*RecorderID).SetValidity(false)
	err = db.CheckID(id)
	if err != ErrInvalidID {
		t.Errorf("wanted %s, got %s", ErrInvalidID, err)
	}

	id.(*RecorderID).SetValidity(true)
	randoError := errors.New("asdf")
	db.Err = randoError
	err = db.CheckID(id)
	if err != randoError {
		t.Errorf("wanted %s, got %s", randoError, err)
	}

	db.Reset()
	if !testStartRecorderDBCleanState(t, db) {
		t.Errorf("TestRecorder reset failed to get to good state")
	}

	if db.Type() != RecorderDBType {
		t.Errorf("RecorderDBType should be %s", RecorderDBType)
	}

	// after all
	db.Reset()
	if !testStartRecorderDBCleanState(t, db) {
		t.Errorf("TestRecorder reset failed to get to good state")
	}
}

func testStartRecorderDBCleanState(t *testing.T, db *RecorderDB) (pass bool) {
	pass = true

	if db.Connected != "" {
		pass = false
		t.Errorf("RecorderDB should start with a \"\" connection")
	}
	if db.ModelUpdates == nil {
		pass = false
		t.Errorf("RecorderDB should start with a ModelUpdates channel")
	}
	if db.Saved == nil || len(db.Saved) != 0 {
		pass = false
		t.Errorf("RecorderDB should start with a 0 Saved array")
	}
	if db.Deleted == nil || len(db.Deleted) != 0 {
		pass = false
		t.Errorf("RecorderDB should start with a 0 Deleted array")
	}
	if db.PopulatedByID == nil || len(db.PopulatedByID) != 0 {
		pass = false
		t.Errorf("RecorderDB should start with a 0 PopulatedByID array")
	}
	if db.PopulatedByField == nil || len(db.PopulatedByField) != 0 {
		pass = false
		t.Errorf("RecorderDB should start with a 0 PopulatedByField array")
	}
	if db.Err != nil {
		pass = false
		t.Errorf("RecorderDB should start with a nil err")
	}

	return
}
