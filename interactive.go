package data

import (
	"sync"
)

var InteractiveType DBType = "interactive"

type InteractiveModel interface {
	Reload() error
}

type InteractiveStore struct {
	Access

	m       sync.Mutex
	objects map[InteractiveModel]bool
}

func NewInteractiveStore(a Access) *InteractiveStore {
	return &InteractiveStore{
		Access:  a,
		objects: make(map[InteractiveModel]bool),
	}
}

func (s *InteractiveStore) Type() DBType {
	return InteractiveType
}

func (s *InteractiveStore) ReloadObjects() {
	s.m.Lock()
	defer s.m.Unlock()

	for object := range s.objects {
		object.Reload()
	}
}

func (s *InteractiveStore) RegisterObject(o InteractiveModel) {
	s.m.Lock()
	defer s.m.Unlock()

	s.objects[o] = true
}
