package data

import (
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
	// implement
}
