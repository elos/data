package mongo_test

import (
	"testing"

	"github.com/elos/data"
	"github.com/elos/data/builtin/mongo"
	"github.com/elos/testing/expect"
)

func TestQuery(t *testing.T) {
	db, err := mongo.New(&mongo.Opts{Addr: "0.0.0.0"})
	expect.NoError("creating db", err, t)
	db.RegisterKind(UserKind, "users")

	u := &User{}
	u.SetID(db.NewID())
	u.Name = "hello"
	if err := db.Save(u); err != nil {
		t.Errorf("Save errors: %s", err)
	}
	t.Logf("Saved user: %v", u)

	// ensure it is in the database
	w := &User{
		Id: u.Id,
	}
	if err := db.PopulateByID(w); err != nil {
		t.Fatalf("PopulateByID error: %v", err)
	}

	if got, want := w.Name, "hello"; got != want {
		t.Fatalf("w.Name: got %q, want %q", got, want)
	}

	iter, err := db.Query(u.Kind()).Select(data.AttrMap{
		"name": "hello",
	}).Execute()
	if got, want := err, error(nil); got != want {
		t.Fatalf("db.Query error: got %v, want %v", got, want)
	}

	u = &User{}
	if ok := iter.Next(u); !ok {
		t.Fatal("iter.Next: got false, want true")
	}

	t.Logf("Retrieved user: %u", u)

	if got, want := u.Name, "hello"; got != want {
		t.Fatalf("u.Name: got %q, want %q", got, want)
	}
}
