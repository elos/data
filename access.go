package data

import "sync"

type Access interface {
	Identifiable() Identifiable
	Store() Store

	SetIdentifiable(Identifiable)
	SetStore(Store)
}

type access struct {
	i Identifiable
	s Store

	*sync.Mutex
}

func NewAccess(i Identifiable, s Store) Access {
	return &access{
		i:     i,
		s:     s,
		Mutex: new(sync.Mutex),
	}
}

func (a *access) SetIdentifiable(i Identifiable) {
	a.Lock()
	defer a.Unlock()

	a.i = i
}

func (a *access) Identifiable() Identifiable {
	a.Lock()
	defer a.Unlock()

	return a.i
}

func (a *access) SetStore(s Store) {
	a.Lock()
	defer a.Unlock()

	a.s = s
}

func (a *access) Store() Store {
	a.Lock()
	defer a.Unlock()

	return a.s
}
