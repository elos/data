package data

import (
	"sync"
)

type LinkKind string
type LinkName string

const (
	MulLink LinkKind = "MANY"
	OneLink LinkKind = "ONE"
)

type Link struct {
	Name    LinkName
	Kind    LinkKind
	Other   Kind
	Inverse LinkName
}

type RelationshipMap map[Kind]map[LinkName]Link

/*
func (s *RelationshipMap) valid() bool {
		for outerKind, links := range *s {
			for innerKind, _ /*linkKind := range links {
				innerLinks, ok := (*s)[innerKind]
				if !ok {
					return false
				}

				_ /*matchingLink, ok = innerLinks[outerKind]

				if !ok {
					return false
				}
			}
		}

	// We used to do complex edge checking, to see that every relation had
	// one which was complimentary, that is now considered over-determined
	// There may be other thinks to check in the future, but for now the
	// Relationship map data structure is self-validating

	return true
}
*/

type versionedRelationshipMap struct {
	*RelationshipMap
	version int

	sync.Mutex
}

func NewSchema(sm *RelationshipMap, version int) (Schema, error) {
	s := &versionedRelationshipMap{
		RelationshipMap: sm,
		version:         version,
	}

	/*
		if !s.valid() {
			return nil, ErrInvalidSchema
		}
	*/

	return s, nil
}

func (s *versionedRelationshipMap) Version() int {
	return s.version
}
