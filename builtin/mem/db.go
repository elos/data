// rudimentary memory db prototype
package mem

import (
	"bytes"
	"fmt"
	"sort"
	"strconv"
	"sync"
	"time"

	"github.com/elos/data"
	"github.com/elos/data/transfer"
	"golang.org/x/net/context"
)

func NewDB() data.DB {
	return &MemDB{
		ChangeHub: data.NewChangeHub(context.TODO()),
		currentID: 0,
		tables:    make(map[data.Kind]map[data.ID]data.Record),
	}
}

func WithData(seed map[data.Kind][]data.Record) data.DB {
	tables := make(map[data.Kind]map[data.ID]data.Record)

	var maxID int64 = 0

	for kind := range seed {
		tables[kind] = make(map[data.ID]data.Record)
	}

	for kind, records := range seed {
		for _, record := range records {
			tables[kind][record.ID()] = record

			id, err := strconv.ParseInt(record.ID().String(), 10, 64)
			if err != nil {
				panic(fmt.Sprintf("parsing id: %v", err))
			}

			if id > maxID {
				maxID = id
			}
		}
	}

	return &MemDB{
		ChangeHub: data.NewChangeHub(context.TODO()),
		currentID: int(maxID),
		tables:    tables,
	}
}

type MemDB struct {
	*data.ChangeHub

	currentID int
	tables    map[data.Kind]map[data.ID]data.Record
}

func (db *MemDB) String() string {
	b := new(bytes.Buffer)
	for k, table := range db.tables {
		fmt.Fprintf(b, "%s:\n", k)
		for _, rec := range table {
			fmt.Fprintf(b, "\t%v\n", rec)
		}
	}
	return b.String()
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
	/*
		var created bool

		if string(r.ID()) == "" {
			r.SetID(db.NewID())
			created = true
		}

		table[r.ID()] = r

		if created {
			db.ChangeHub.Notify(data.NewCreate(r))
		} else {
			db.ChangeHub.Notify(data.NewUpdate(r))
		}
	*/

	if string(r.ID()) == "" {
		return data.ErrInvalidID
	}

	_, existed := table[r.ID()]

	table[r.ID()] = r

	if existed {
		db.ChangeHub.Notify(data.NewUpdate(r))
	} else {
		db.ChangeHub.Notify(data.NewCreate(r))
	}

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
	db.ChangeHub.Notify(data.NewDelete(r))

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
	kind               data.Kind
	db                 *MemDB
	wheres             map[string]interface{}
	limit, skip, batch int
	order              []string
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

	buffer := make(chan data.Record)

	// buffer and simulate skipping/limitting

	go func() {
		index := -1 // so that it starts at 0 on the first receive
		for r := range out {
			index++

			if q.skip != 0 && index < q.skip {
				continue // don't forward
			}

			if q.limit != 0 && index >= q.limit {
				close(buffer)
				return
			}

			buffer <- r
		}
		close(buffer)
	}()

	return Iter(sorted(buffer, q.order...)), nil
}

type rMap struct {
	m map[string]interface{}
	data.Record
}

type byFields struct {
	records []*rMap
	fields  []string
}

func (b *byFields) Len() int { return len(b.records) }

func (b *byFields) Less(i, j int) bool {
	for _, f := range b.fields {
		if lessThan(b.records[i].m[f], b.records[j].m[f]) {
			return true
		}

		if greaterThan(b.records[i].m[f], b.records[j].m[f]) {
			return false
		}
	}

	return true
}

func (b *byFields) Swap(i, j int) {
	b.records[i], b.records[j] = b.records[j], b.records[i]
}

func lessThan(i, j interface{}) bool {
	switch i.(type) {
	case bool:
		return false
	case int:
		return i.(int) < j.(int)
	case float64:
		return i.(float64) < j.(float64)
	case string:
		return i.(string) < j.(string)
	case time.Time:
		return i.(time.Time).Before(j.(time.Time))
	default:
		return false
	}
}

func greaterThan(i, j interface{}) bool {
	switch i.(type) {
	case bool:
		return false
	case int:
		return i.(int) > j.(int)
	case float64:
		return i.(float64) > j.(float64)
	case string:
		return i.(string) > j.(string)
	case time.Time:
		return i.(time.Time).After(j.(time.Time))
	default:
		return false
	}
}

func sorted(in <-chan data.Record, fields ...string) <-chan data.Record {
	if len(fields) == 0 {
		return in
	}

	out := make(chan data.Record)

	go func() {
		records := make([]*rMap, 0)

		for r := range in {
			m := &rMap{
				m:      make(map[string]interface{}),
				Record: r,
			}
			transfer.TransferAttrs(r, &m.m)
			records = append(records, m)
		}

		b := &byFields{
			records: records,
			fields:  fields,
		}

		sort.Sort(b)

		for _, r := range records {
			out <- r.Record
		}
	}()
	return out
}

func contains(r data.Record, field string, v interface{}) bool {
	m := make(map[string]interface{})
	if err := transfer.TransferAttrs(r, &m); err != nil {
		panic(fmt.Sprintf("trying to transfer from %+v of type %T error: %v", r, r, err))
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
	q.skip = i
	return q
}

func (q *memQuery) Limit(i int) data.Query {
	q.limit = i
	return q
}

func (q *memQuery) Batch(i int) data.Query {
	q.batch = i
	return q
}

func (q *memQuery) Select(m data.AttrMap) data.Query {
	for k, v := range m {
		q.wheres[k] = v
	}
	return q
}

func (q *memQuery) Order(fields ...string) data.Query {
	q.order = fields
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

func Slice(i data.Iterator, constructor func() data.Record) []data.Record {
	results := make([]data.Record, 0)

	r := constructor()
	for i.Next(r) {
		results = append(results, r)
		r = constructor()
	}

	return results
}
