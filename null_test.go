package data

import "testing"

func TestNullID(t *testing.T) {
	if NewNullID("").Valid() != true {
		t.Errorf("NullID should always be valid")
	}

	_ = func() ID { return NewNullID("") }
}

func TetNullDB(t *testing.T) {
	_ = func() DB { return NewNullDB() }
}

func TestNullSchema(t *testing.T) {
	_ = func() Schema { return NewNullSchema() }
}

func TestNullStore(t *testing.T) {
	_ = func() Store { return NewNullStore() }
}
