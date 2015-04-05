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

type Access interface {
	Store
	Client() Client
}

func NewAnonAccess(s Store) Access {
	return NewAccess(AnonClient, s)
}

type access struct {
	client Client
	store  Store

	*sync.Mutex
}

func NewAccess(c Client, s Store) Access {
	return &access{
		client: c,
		store:  s,
		Mutex:  new(sync.Mutex),
	}
}

func (a *access) Client() Client {
	return a.client
}

func (a *access) Save(m Model) error {
	if m.CanWrite(a.client) {
		return a.store.Save(m)
	} else {
		return ErrAccessDenial
	}
}

func (a *access) Delete(m Model) error {
	if m.CanWrite(a.client) {
		return a.store.Delete(m)
	} else {
		return ErrAccessDenial
	}
}

func (a *access) PopulateByID(m Model) error {
	// todo optimize
	temp, err := a.store.ModelFor(m.Kind())
	if err != nil {
		return err
	}

	temp.SetID(m.ID())
	if err = a.store.PopulateByID(temp); err != nil {
		return err
	}

	if temp.CanRead(a.client) {
		return a.store.PopulateByID(m)
	} else {
		return ErrAccessDenial
	}
}

func (a *access) PopulateByField(s string, v interface{}, m Model) error {
	temp, err := a.store.ModelFor(m.Kind())
	if err != nil {
		return err
	}

	if err = a.store.PopulateByField(s, v, temp); err != nil {
		return err
	}

	if temp.CanRead(a.client) {
		return a.store.PopulateByField(s, v, m)
	} else {
		return ErrAccessDenial
	}
}

func (a *access) RegisterForChanges(c Client) *chan *Change {
	return a.store.RegisterForChanges(c)
}

func (a *access) Register(k Kind, c ModelConstructor) {
	a.store.Register(k, c)
}

func (a *access) Unmarshal(k Kind, attrs AttrMap) (Model, error) {
	return a.store.Unmarshal(k, attrs)
}

func (a *access) ModelFor(k Kind) (Model, error) {
	return a.store.ModelFor(k)
}

func (a *access) NewID() ID {
	return a.store.NewID()
}

func (a *access) ParseID(s string) (ID, error) {
	return a.store.ParseID(s)
}

func (a *access) Query(k Kind) ModelQuery {
	return a.store.Query(k)
}

func (a *access) Registered() []Kind {
	return a.store.Registered()
}

func (a *access) Type() DBType {
	return a.store.Type()
}
