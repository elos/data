package data

import "testing"

func TestCompatible(t *testing.T) {
	if !Compatible(NewNullModel(), NewNullModel()) {
		t.Errorf("ExampleModels should be compatible")
	}
}

var name LinkName = "test"

var r = RelationshipMap{
	NullKind: {
		name: Link{
			Name:  name,
			Kind:  MulLink,
			Other: NullKind,
		},
	},
}

var rInverse = RelationshipMap{
	NullKind: {
		name: Link{
			Name:    name,
			Kind:    MulLink,
			Other:   NullKind,
			Inverse: name,
		},
	},
}

var blankR = RelationshipMap{}

var badR = RelationshipMap{
	NullKind: {
		name: Link{
			Name:  name,
			Kind:  MulLink,
			Other: "trash",
		},
	},
}

var (
	e1 = NewNullModel()
	e2 = NewNullModel()
)

func TestPossibleLink(t *testing.T) {
	ok, err := possibleLink(&blankR, e1, "asdf")
	if ok != false {
		t.Errorf("should not be a possibleLink")
	}

	if _, ok = err.(UndefinedKindError); !ok {
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

func TestLinkFor(t *testing.T) {
	_, err := (&blankR).linkFor(e1, name)
	if _, ok := err.(UndefinedKindError); !ok {
		t.Errorf("should error")
	}

	l, err := (&r).linkFor(e1, name)
	if err != nil {
		t.Errorf("err should be nil")
	}
	if l != r[NullKind][name] {
		t.Errorf("linkfor did not retrieved the correct link")
	}
	// Not that multiple links of the same kind need not be tested as the map
	// data structure explicitly disallows this.
}

func TestLink(t *testing.T) {
	temp := exampleLink
	defer func(t func(m Model, l Link) error) {
		exampleLink = t
	}(temp)

	m := make(chan Model, 2)
	l := make(chan Link)

	exampleLink = func(model Model, link Link) error {
		go func() { m <- model }()
		go func() { l <- link }()
		return nil
	}

	// Correct call
	err := (&rInverse).Link(e1, e2, name)
	if err != nil {
		t.Errorf("err should be nil")
	}

	pe2 := <-m // note calls link of m1 first (a bit implementation-specific)
	pe1 := <-m
	if e2 != pe2.(*NM) || e1 != pe1.(*NM) {
		t.Errorf("Link should call Link on the model with the other model")
	}
	pl1 := <-l
	pl2 := <-l
	correctLink := rInverse[NullKind][name]
	if pl1 != correctLink || pl2 != correctLink {
		t.Errorf("Link didn't pass correct link")
	}
}

func TestUnlink(t *testing.T) {
	temp := exampleUnlink
	defer func(t func(m Model, l Link) error) {
		exampleUnlink = t
	}(temp)

	m := make(chan Model, 2)
	l := make(chan Link)

	exampleUnlink = func(model Model, link Link) error {
		go func() { m <- model }()
		go func() { l <- link }()
		return nil
	}

	// Correct call
	err := (&rInverse).Unlink(e1, e2, name)
	if err != nil {
		t.Errorf("err should be nil")
	}

	pe2 := <-m // note calls link of m1 first (a bit implementation-specific)
	pe1 := <-m
	if e2 != pe2.(*NM) || e1 != pe1.(*NM) {
		t.Errorf("Link should call Link on the model with the other model")
	}
	pl1 := <-l
	pl2 := <-l
	correctLink := rInverse[NullKind][name]
	if pl1 != correctLink || pl2 != correctLink {
		t.Errorf("Link didn't pass correct link")
	}
}
