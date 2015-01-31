package data

import "sync"

type LinkKind string

const (
	MulLink LinkKind = "MANY"
	OneLink LinkKind = "ONE"
)

type RelationshipMap map[Kind]map[Kind]LinkKind

func (s *RelationshipMap) valid() bool {
	for outerKind, links := range *s {
		for innerKind, _ /*linkKind*/ := range links {
			innerLinks, ok := (*s)[innerKind]
			if !ok {
				return false
			}

			_ /*matchingLink*/, ok = innerLinks[outerKind]

			if !ok {
				return false
			}
		}
	}

	return true
}

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

	if !s.valid() {
		return nil, ErrInvalidSchema
	}

	return s, nil
}

func (s *versionedRelationshipMap) Version() int {
	return s.version
}
