package data

import (
	"testing"
)

func TestSchema(t *testing.T) {
	r := RelationshipMap{
		ExampleKind: {
			"name": Link{
				Name:    "link",
				Kind:    MulLink,
				Other:   ExampleKind,
				Inverse: "name",
			},
		},
	}

	s, err := NewSchema(&r, 1)

	if s == nil || err != nil {
		t.Errorf("NewSchema should always work") // note there is no validation of the
		// relationship map at thie point
	}

	if s.Version() != 1 {
		t.Errorf("Schema does not properly store version")
	}

}
