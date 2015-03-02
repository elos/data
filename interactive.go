package data

import (
	"sync"
)

var SpaceType DBType = "space"

type InteractiveModel interface {
	Model
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
	return SpaceType
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
