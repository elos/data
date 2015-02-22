package data

import (
	"encoding/json"
	"sync"
)

type Store interface {
	// DB Info
	Type() DBType

	// Schema Model Management
	Register(Kind, ModelConstructor)
	ModelFor(Kind) (Model, error)
	Unmarshal(Kind, AttrMap) (Model, error)
	Registered() []Kind

	// Model Persistence
	Save(Model) error
	Delete(Model) error
	PopulateByID(Model) error
	PopulateByField(string, interface{}, Model) error
	RegisterForChanges(Client) *chan *Change

	NewID() ID
	Query(Kind) ModelQuery
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

func (s *store) Registered() []Kind {
	s.Lock()
	defer s.Unlock()
	k := make([]Kind, 0)
	for kind, _ := range s.registered {
		k = append(k, kind)
	}

	return k
}

func (s *store) Save(m Model) error {
	return s.DB.Save(m)
}

func (s *store) Delete(m Model) error {
	return s.DB.Delete(m)
}

func (s *store) PopulateByID(m Model) error {
	return s.DB.PopulateByID(m)
}

func (s *store) PopulateByField(str string, v interface{}, m Model) error {
	return s.DB.PopulateByField(str, v, m)
}

func (s *store) NewID() ID {
	return s.DB.NewID()
}

type qb struct {
	RecordQuery
}

func (q *qb) Execute() (ModelIterator, error) {
	i, err := q.RecordQuery.Execute()
	return &ib{RecordIterator: i}, err
}

func (q *qb) Select(attrs AttrMap) ModelQuery {
	q.RecordQuery.Select(attrs)
	return q
}

func (q *qb) Limit(i int) ModelQuery {
	q.RecordQuery.Limit(i)
	return q
}

func (q *qb) Skip(i int) ModelQuery {
	q.RecordQuery.Skip(i)
	return q
}

func (q *qb) Batch(i int) ModelQuery {
	q.RecordQuery.Batch(i)
	return q
}

type ib struct {
	RecordIterator
}

func (i *ib) Next(m Model) bool {
	return i.RecordIterator.Next(m)
}

func (i *ib) Close() error {
	return i.RecordIterator.Close()
}

func (s *store) Query(k Kind) ModelQuery {
	return &qb{RecordQuery: s.DB.NewQuery(k)}
}
