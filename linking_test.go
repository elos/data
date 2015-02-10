package data

import (
	"testing"
)

func TestCompatible(t *testing.T) {
	if !Compatible(NewExampleModel(), NewExampleModel()) {
		t.Errorf("ExampleModels should be compatible")
	}
}

var name LinkName = "test"

var r = RelationshipMap{
	ExampleKind: {
		name: Link{
			Name:  name,
			Kind:  MulLink,
			Other: ExampleKind,
		},
	},
}

var blankR = RelationshipMap{}

var badR = RelationshipMap{
	ExampleKind: {
		name: Link{
			Name:  name,
			Kind:  MulLink,
			Other: "trash",
		},
	},
}

var e1 = NewExampleModel()

// e2 := NewExampleModel()

func TestPossibleLink(t *testing.T) {
	if r.valid() != true || blankR.valid() != true || badR.valid() != true {
		t.Errorf("all relationship maps should be valid")
	}

	ok, err := possibleLink(&blankR, e1, "asdf")
	if ok != false {
		t.Errorf("should not be a possibleLink")
	}
	if err != ErrUndefinedKind {
		t.Errorf("possible link shouldn't have found a link in the map")
	}

	ok, err = possibleLink(&badR, e1, "asdfa")
	if ok != false {
		t.Errorf("should not be possibleLink")
	}
	if err != ErrUndefinedLink {
		t.Errorf("should be undefined link cause wrong name")
	}

	ok, err = possibleLink(&r, e1, name)
	if ok != true || err != nil {
		t.Errorf("Should be a link")
	}
}
