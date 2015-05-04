package mongo_test

import (
	"testing"
	"time"

	"github.com/elos/data"
	"github.com/elos/data/builtin/mongo"
	"github.com/elos/testing/expect"
)

func TestSave(t *testing.T) {
	db, err := mongo.New(&mongo.Opts{Addr: "localhost"})
	expect.NoError("creating db", err, t)
	db.RegisterKind(UserKind, "users")

	testString := "asdfasdfASfaDF2309581934"
	testTime := time.Now()

	u := &User{}
	u.Name = testString
	u.CreatedAt = testTime

	if err := db.Save(u); err != data.ErrInvalidID {
		t.Errorf("Save errors: %s", err.Error())
	}

	if err := db.Save(u); err != data.ErrInvalidID {
		t.Errorf("Mongo should not choke on bad ids")
	}

	u.SetID(db.NewID())

	if err := db.Save(u); err != nil {
		t.Errorf("Model failed to save, err: %s", err.Error())
	}

	// retrieval
	r := &User{}
	r.SetID(u.ID())

	if err := db.PopulateByID(r); err != nil {
		t.Errorf("DB failed to populate model by id, %s", err.Error())
	}

	if r.Name != testString {
		t.Errorf("Error with data, wanted %s, got %s", testString, r.Name)
	}

	if !(r.CreatedAt.Sub(testTime) < 1*time.Second) {
		t.Errorf("Error with data, wanted %d, got %d", testTime, r.CreatedAt)
	}

}

func TestDelete(t *testing.T) {
	db, err := mongo.New(&mongo.Opts{Addr: "localhost"})
	expect.NoError("creating db", err, t)
	db.RegisterKind(UserKind, "users")

	testString := "aksjdkjasfd"
	u := &User{}
	u.Name = testString

	if err := db.Delete(u); err != data.ErrInvalidID {
		t.Errorf("Delete should recognize a bad id")
	}

	u.SetID(db.NewID())
	if err := db.Save(u); err != nil {
		t.Errorf("Save should go off fine")
	}

	// can assume it is there, testing in save_test

	if err := db.Delete(u); err != nil {
		t.Errorf("Delete should work, but errored: %s", err.Error())
	}

	r := &User{}
	r.SetID(u.ID())
	if err := db.PopulateByID(r); err != data.ErrNotFound {
		t.Errorf("the delete should have removed the model")
	}
}

func TestPopulate(t *testing.T) {
	db, err := mongo.New(&mongo.Opts{Addr: "localhost"})
	expect.NoError("creating db", err, t)
	db.RegisterKind(UserKind, "users")

	testString := "aklksadjf234sjdf"

	u := &User{}
	u.Name = testString
	id := db.NewID()
	u.SetID(id)

	err = db.Save(u)
	expect.NoError("saving model", err, t)

	u = &User{}

	if err := db.PopulateByID(u); err != data.ErrInvalidID {
		t.Errorf("PopulateByID should reject a model with an invalid ID")
	}

	u.SetID(id)

	if err := db.PopulateByID(u); err != nil {
		t.Errorf("PopulateByID should return nil on a valid model, but got %s", err.Error())
	}

	if u.Name != testString {
		t.Errorf("PopulateByID failed to correctly populate model, got %s, wanted: %s", u.Name, testString)
	}

	u = &User{}
	if err := db.PopulateByField("name", testString, u); err != nil {
		t.Errorf("PopulateByField should return nil on valid model, but go %s", err.Error)
	}

	if u.Name != testString {
		t.Errorf("Expected %s, got %s", testString, u.Name)
	}

	if u.ID().String() != id.String() {
		t.Errorf("ID of the populated model didn't match, got: %s, wanted: %s", u.ID().String(), id.String())
	}
}
