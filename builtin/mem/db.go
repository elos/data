// rudimentary memory db prototype
package mem

import (
	"fmt"
	"log"
	"sync"

	"github.com/elos/data"
	"github.com/elos/data/transfer"
)

func NewDB() data.DB {
	return &MemDB{
		ChangePub: data.NewChangePub(),
		currentID: 0,
		tables:    make(map[data.Kind]map[data.ID]data.Record),
	}
}

type MemDB struct {
	*data.ChangePub

	currentID int
	tables    map[data.Kind]map[data.ID]data.Record
}

func (db *MemDB) NewID() data.ID {
	db.currentID += 1
	return data.ID(fmt.Sprintf("%d", db.currentID))
}

func (db *MemDB) ParseID(id string) (data.ID, error) {
	return data.ID(id), nil
}

func (db *MemDB) Save(r data.Record) error {
	table, ok := db.tables[r.Kind()]
	if !ok {
		table = make(map[data.ID]data.Record)
		db.tables[r.Kind()] = table
	}

	table[r.ID()] = r

	db.ChangePub.Notify(data.NewUpdate(r))

	return nil
}

func (db *MemDB) Delete(r data.Record) error {
	table, ok := db.tables[r.Kind()]
	if !ok {
		return nil
	}

	_, ok = table[r.ID()]
	if !ok {
		return nil
	}

	delete(table, r.ID())
	db.ChangePub.Notify(data.NewDelete(r))

	return nil
}

func (db *MemDB) PopulateByID(r data.Record) error {
	table, ok := db.tables[r.Kind()]
	if !ok {
		return data.ErrNotFound
	}

	storedRecord, ok := table[r.ID()]
	if !ok {
		return data.ErrNotFound
	}

	return transfer.TransferAttrs(storedRecord, r)
}

func (db *MemDB) PopulateByField(field string, v interface{}, r data.Record) error {
	table, ok := db.tables[r.Kind()]
	if !ok {
		return data.ErrNotFound
	}

	for _, stored := range table {
		if contains(stored, field, v) {
			return transfer.TransferAttrs(stored, r)
		}
	}

	return data.ErrNotFound
}

func (db *MemDB) Query(k data.Kind) data.Query {
	return newMemQuery(k, db)
}

func newMemQuery(k data.Kind, db *MemDB) *memQuery {
	return &memQuery{
		kind:   k,
		db:     db,
		wheres: make(map[string]interface{}),
	}
}

type memQuery struct {
	kind   data.Kind
	db     *MemDB
	wheres map[string]interface{}
}

func (q *memQuery) Execute() (data.Iterator, error) {
	return q.exec()
}

func (q *memQuery) exec() (data.Iterator, error) {
	table, ok := q.db.tables[q.kind]
	if !ok {
		out := make(chan data.Record)
		defer close(out)
		return Iter(out), nil
	}

	in := make(chan data.Record, len(table))

	for _, r := range table {
		in <- r
	}
	defer close(in)

	var out <-chan data.Record = in

	for s, v := range q.wheres {
		out = filter(out, s, v)
	}

	return Iter(out), nil
}

func contains(r data.Record, field string, v interface{}) bool {
	m := make(map[string]interface{})
	if err := transfer.TransferAttrs(r, &m); err != nil {
		log.Fatal(err)
	}
	return equals(m[field], v)
}

func equals(v interface{}, w interface{}) bool {
	switch v.(type) {
	case int:
		return v == w
	case float64:
		return v == w
	case bool:
		return v == w
	case string:
		return v == w
	case []interface{}:
		if len(v.([]interface{})) != len(w.([]interface{})) {
			return false
		}

		for i := range v.([]interface{}) {
			if !equals(v.([]interface{})[i], w.([]interface{})[i]) {
				return false
			}
		}

		return true
	default:
		return v == w
	}
}

// this is slow af
func filter(in <-chan data.Record, field string, v interface{}) <-chan data.Record {
	out := make(chan data.Record)

	go func() {
		for r := range in {
			if contains(r, field, v) {
				out <- r
			}
		}

		close(out)
	}()

	return out
}

func (q *memQuery) Skip(i int) data.Query {
	return q
}

func (q *memQuery) Limit(i int) data.Query {
	return q
}

func (q *memQuery) Batch(i int) data.Query {
	return q
}

func (q *memQuery) Select(m data.AttrMap) data.Query {
	for k, v := range m {
		q.wheres[k] = v
	}
	return q
}

func Iter(c <-chan data.Record) data.Iterator {
	return &memIter{
		inbound: c,
	}
}

type memIter struct {
	inbound <-chan data.Record
	sync.Mutex
}

func (i *memIter) Next(r data.Record) bool {
	in, ok := <-i.inbound

	if ok {
		transfer.TransferAttrs(in, r)
	}

	return ok
}

func (i *memIter) Close() error {
	return nil
}
