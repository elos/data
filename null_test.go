package data

import "testing"

func TestNullID(t *testing.T) {
	if NewNullID("").Valid() != true {
		t.Errorf("NullID should always be valid")
	}

	_ = func() ID { return NewNullID("") }
}

func TestNullDB(t *testing.T) {
	_ = func() DB { return NewNullDB() }

	db := NewNullDB()

	if db.Type() != nullDBType {
		t.Errorf("Type() failed")
	}

	if db.Connect("random") != nil {
		t.Errorf("Connect should return nil")
	}

	updates := db.RegisterForUpdates(NewExampleModel())

	if updates == nil {
		t.Errorf("RegisterForUpdates should return a real non-null channel")
	}

	id := db.NewID()
	id, ok := id.(nullID)
	if !ok || id == nil {
		t.Errorf("NewID should return a nullId")
	}

	if err := db.CheckID(id); err != nil {
		t.Errorf("NullDB should always return a non-nil error")
	}

	r := NewExampleModel()

	if err := db.Save(r); err != nil {
		t.Errorf("Save should always return a non-nil error")
	}

	if err := db.Delete(r); err != nil {
		t.Errorf("Delete should always return a non-nil error")
	}

	if err := db.PopulateByID(r); err != nil {
		t.Errorf("PopulateByID should always return a non-nil error")
	}

	if err := db.PopulateByField("ads", "afd", r); err != nil {
		t.Errorf("PopulateByField should always return a non-nil error")
	}

	if query := db.NewQuery(ExampleKind); query != nil {
		t.Errorf("Query should always return nil")
	}

}

func TestNullSchema(t *testing.T) {
	_ = func() Schema { return NewNullSchema() }
	// just wraps a real schema
}

func TestNullStore(t *testing.T) {
	_ = func() Store { return NewNullStore() }

	s := NewNullStore()

	if s == nil {
		t.Errorf("NewNullStore should not return nil")
	}

	if s.Type() != nullDBType {
		t.Errorf("Type should be nullDBType")
	}

	s.Register(ExampleKind, NewEM)

	one, two := s.ModelFor(ExampleKind)

	if one != nil || two != nil {
		t.Errorf("ModelFor should return two nils")
	}

	one, two = s.Unmarshal(ExampleKind, AttrMap{})

	if one != nil || two != nil {
		t.Errorf("Unmarshal should return two nils")
	}

	s = NewNullStoreWithType("test")
	if s.Type() != "test" {
		t.Errorf("NewNullStoreWithType should have custom type")
	}
}
