package data

import (
	"encoding/json"
	"sync"
)

type Store interface {
	DB
	Schema

	Register(Kind, ModelConstructor)
	ModelFor(Kind) (Model, error)
	Unmarshal(Kind, AttrMap) (Model, error)
}

type store struct {
	DB
	Schema
	registered map[Kind]ModelConstructor
	sync.Mutex
}

func NewStore(db DB, s Schema) *store {
	return &store{
		DB:         db,
		Schema:     s,
		registered: make(map[Kind]ModelConstructor),
	}
}

type ModelConstructor func(Store) (Model, error)

func (s *store) Register(k Kind, c ModelConstructor) {
	s.Lock()
	defer s.Unlock()

	s.registered[k] = c
}

func (s *store) ModelFor(kind Kind) (Model, error) {
	s.Lock()
	defer s.Unlock()
	c, ok := s.registered[kind]

	if !ok {
		return nil, ErrUndefinedKind
	}

	return c(s)
}

func (s *store) Unmarshal(k Kind, attrs AttrMap) (Model, error) {
	bytes, _ := json.Marshal(attrs)
	m, err := s.ModelFor(k)
	if err != nil {
		return m, err
	}

	return m, json.Unmarshal(bytes, m)
}

func (s *store) RegisteredModels() []Kind {
	s.Lock()
	defer s.Unlock()
	k := make([]Kind, 0)
	for kind, _ := range s.registered {
		k = append(k, kind)
	}

	return k
}
