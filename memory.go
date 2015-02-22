package data

import (
	"sync"
)

var SpaceType DBType = "space"

type MemoryObject interface {
	Reload() error
}

type MemoryStore struct {
	Access

	m       sync.Mutex
	objects map[MemoryObject]bool
}

func NewMemoryStore(a Access) *MemoryStore {
	return &MemoryStore{
		Access:  a,
		objects: make(map[MemoryObject]bool),
	}
}

func (s *MemoryStore) Type() DBType {
	return SpaceType
}

func (s *MemoryStore) ReloadObjects() {
	s.m.Lock()
	defer s.m.Unlock()

	for object := range s.objects {
		object.Reload()
	}
}

func (s *MemoryStore) RegisterObject(o MemoryObject) {
	s.m.Lock()
	defer s.m.Unlock()

	s.objects[o] = true
}
