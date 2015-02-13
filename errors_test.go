package data

import (
	"strings"
	"testing"
)

func TestNewAttrError(t *testing.T) {
	e := NewAttrError("first", "second")
	expected := "attribute first must second"
	actual := e.Error()
	if actual != expected {
		t.Errorf("got: %s, expected: %s", actual, expected)
	}
}

func TestNewLinkError(t *testing.T) {
	m1 := NewNullModel()
	m1.String = "one"
	m2 := NewNullModel()
	m2.String = "two"

	e := NewLinkError(m1, m2, *NullLink)

	if !strings.Contains(e.Error(), "could not be linked") {
		t.Errorf("Something is wrong with LinkError")
	}
}
