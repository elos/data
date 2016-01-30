package mem_test

import (
	"testing"
	"time"

	"github.com/elos/data"
	"github.com/elos/data/builtin/mem"
)

const TestRecordKind data.Kind = "test"

type TestRecord struct {
	Id string

	Name       string
	Count      int
	Percentage float64
	Time       time.Time
}

func (tr *TestRecord) Kind() data.Kind {
	return TestRecordKind
}

func (tr *TestRecord) SetID(id data.ID) {
	tr.Id = id.String()
}

func (tr *TestRecord) ID() data.ID {
	return data.ID(tr.Id)
}

func TestQueryLimit(t *testing.T) {
	db := mem.NewDB()

	tr1 := &TestRecord{Name: "one"}
	tr2 := &TestRecord{Name: "two"}

	tr1.SetID(db.NewID())
	tr2.SetID(db.NewID())

	if err := db.Save(tr1); err != nil {
		t.Fatal(err)
	}

	if err := db.Save(tr2); err != nil {
		t.Fatal(err)
	}

	iter, err := db.Query(TestRecordKind).Limit(1).Execute()
	if err != nil {
		t.Fatal(err)
	}

	tr := &TestRecord{}
	if ok := iter.Next(tr); !ok {
		t.Fatal("Expected there to be at least one record in the query")
	}
	if ok := iter.Next(tr); ok {
		t.Fatal("Expected there to be at most one record in the query")
	}

	if err = iter.Close(); err != nil {
		t.Fatal(err)
	}
}

func TestQuerySkip(t *testing.T) {
	db := mem.NewDB()

	tr1 := &TestRecord{Name: "one"}
	tr2 := &TestRecord{Name: "two"}

	tr1.SetID(db.NewID())
	tr2.SetID(db.NewID())

	if err := db.Save(tr1); err != nil {
		t.Fatal(err)
	}

	if err := db.Save(tr2); err != nil {
		t.Fatal(err)
	}

	iter, err := db.Query(TestRecordKind).Skip(1).Execute()
	if err != nil {
		t.Fatal(err)
	}

	tr := &TestRecord{}
	if ok := iter.Next(tr); !ok {
		t.Fatal("Expected there to be at least one record in the query")
	}
	if ok := iter.Next(tr); ok {
		t.Fatal("Expected there to be at most one record in the query")
	}

	if err = iter.Close(); err != nil {
		t.Fatal(err)
	}
}
