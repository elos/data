package data

import (
	"testing"
)

func TestMap(t *testing.T) {
	e := NewNullModel()
	kmap := Map(e)

	if kmap[NullKind] != e {
		t.Errorf("TestMap failed")
	}
}

func TestNewChange(t *testing.T) {
	e := NewNullModel()
	c := NewChange(Update, e)

	if c.ChangeKind != Update || c.Record != e {
		t.Errorf("NewChange failed")
	}
}
