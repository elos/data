package data

import (
	"testing"
)

func TestMap(t *testing.T) {
	e := NewExampleModel()
	kmap := Map(e)

	if kmap[ExampleKind] != e {
		t.Errorf("TestMap failed")
	}
}

func TestNewChange(t *testing.T) {
	e := NewExampleModel()
	c := NewChange(Update, e)

	if c.ChangeKind != Update || c.Record != e {
		t.Errorf("NewChange failed")
	}
}
