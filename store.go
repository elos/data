package data

type Store interface {
	DB
}

type store struct {
	DB
}

func NewStore(db DB) *store {
	return &store{
		DB: db,
	}
}

func (s *store) Query(k Kind) ModelQuery {
	return &qb{RecordQuery: s.DB.NewQuery(k)}
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
