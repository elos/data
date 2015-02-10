package data

import "sync"

const Anonymous Kind = "anonymous"

type anonClient NullID

func (a anonClient) ID() ID {
	return NullID(a)
}

func (a anonClient) Kind() Kind {
	return Anonymous
}

// only one anon client;
var AnonClient = anonClient(NewNullID("anon"))

func NewAnonAccess(s Store) *Access {
	return NewAccess(AnonClient, s)
}

type Access struct {
	Client
	Store

	*sync.Mutex
}

func NewAccess(c Client, s Store) *Access {
	return &Access{
		Client: c,
		Store:  s,
		Mutex:  new(sync.Mutex),
	}
}

func (a *Access) Save(m Model) error {
	if m.CanWrite(a.Client) {
		return a.Store.Save(m)
	} else {
		return ErrAccessDenial
	}
}

func (a *Access) Delete(m Model) error {
	if m.CanWrite(a.Client) {
		return a.Store.Delete(m)
	} else {
		return ErrAccessDenial
	}
}

func (a *Access) PopulateByID(m Model) error {
	// todo optimize
	temp, err := a.Store.ModelFor(m.Kind())
	if err != nil {
		return err
	}

	temp.SetID(m.ID())
	if err = a.Store.PopulateByID(temp); err != nil {
		return err
	}

	if temp.CanRead(a.Client) {
		return a.Store.PopulateByID(m)
	} else {
		return ErrAccessDenial
	}
}

func (a *Access) PopulateByField(s string, v interface{}, m Model) error {
	temp, err := a.Store.ModelFor(m.Kind())
	if err != nil {
		return err
	}

	if err = a.Store.PopulateByField(s, v, temp); err != nil {
		return err
	}

	if temp.CanRead(a.Client) {
		return a.Store.PopulateByField(s, v, m)
	} else {
		return ErrAccessDenial
	}
}

func (a *Access) RegisterForUpdates(Identifiable) *chan *Change {
	return a.Store.RegisterForUpdates(a.Client)
}

func (a *Access) Unmarshal(k Kind, attrs AttrMap) (Model, error) {
	return a.Store.Unmarshal(k, attrs)
}
