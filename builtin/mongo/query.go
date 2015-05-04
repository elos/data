package mongo

import (
	"sync"

	"github.com/elos/data"
	"gopkg.in/mgo.v2"
)

type Query struct {
	db                 *DB
	kind               data.Kind
	match              data.AttrMap
	limit, skip, batch int
	m                  sync.Mutex
}

func (q *Query) Execute() (data.Iterator, error) {
	q.m.Lock()
	defer q.m.Unlock()

	s, err := q.db.Fork()
	if err != nil {
		return nil, err
	}
	defer s.Close()

	c, err := q.db.Collection(s, q.kind)
	if err != nil {
		return nil, err
	}

	mgoQuery := c.Find(q.match)

	if q.limit != 0 {
		mgoQuery.Limit(q.limit)
	}

	if q.skip != 0 {
		mgoQuery.Skip(q.skip)
	}

	if q.batch != 0 {
		mgoQuery.Batch(q.batch)
	}

	return newIter(mgoQuery.Iter()), nil
}

func (q *Query) Select(am data.AttrMap) data.Query {
	q.m.Lock()
	defer q.m.Unlock()

	q.match = am
	return q
}

func (q *Query) Limit(i int) data.Query {
	q.m.Lock()
	defer q.m.Unlock()

	q.limit = i
	return q
}

func (q *Query) Skip(i int) data.Query {
	q.m.Lock()
	defer q.m.Unlock()

	q.skip = i
	return q
}

func (q *Query) Batch(i int) data.Query {
	q.m.Lock()
	defer q.m.Unlock()

	q.batch = i
	return q
}

type iter struct {
	iter *mgo.Iter
	sync.Mutex
}

func newIter(i *mgo.Iter) data.Iterator {
	return &iter{iter: i}
}

func (i *iter) Next(r data.Record) bool {
	return i.iter.Next(r)
}

func (i *iter) Close() error {
	return i.iter.Close()
}
