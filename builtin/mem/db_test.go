package mem_test

import (
	"testing"
	"time"

	"github.com/elos/data"
	"github.com/elos/data/builtin/mem"
)

const TestRecordKind data.Kind = "test"

type s struct {
	Foo string
}

type TestRecord struct {
	Id string

	Name       string
	Count      int
	Percentage float64
	Time       time.Time
	True       bool
	Ptr        *s
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

func TestQueryOrder(t *testing.T) {
	db := mem.WithData(map[data.Kind][]data.Record{
		TestRecordKind: []data.Record{
			&TestRecord{
				Id:    "1",
				Name:  "third",
				Count: 3,
			},
			&TestRecord{
				Id:    "2",
				Name:  "second",
				Count: 2,
			},
			&TestRecord{
				Id:    "3",
				Name:  "first",
				Count: 1,
			},
		},
	})

	iter, err := db.Query(TestRecordKind).Order("Count").Execute()
	if err != nil {
		t.Fatalf("db.Query error: %v", err)
	}

	record := new(TestRecord)

	if got, want := iter.Next(record), true; got != want {
		t.Fatalf("iter.Next: got %t, want %t", got, want)
	}

	if got, want := record.Name, "first"; got != want {
		t.Errorf("record.Name: got %q, want %q", got, want)
	}

	if got, want := iter.Next(record), true; got != want {
		t.Fatalf("iter.Next: got %t, want %t", got, want)
	}

	if got, want := record.Name, "second"; got != want {
		t.Errorf("record.Name: got %q, want %q", got, want)
	}

	if got, want := iter.Next(record), true; got != want {
		t.Fatalf("iter.Next: got %t, want %t", got, want)
	}

	if got, want := record.Name, "third"; got != want {
		t.Errorf("record.Name: got %q, want %q", got, want)
	}

	if got, want := iter.Next(record), false; got != want {
		t.Fatalf("iter.Next: got %t, want %t", got, want)
	}

	if err := iter.Close(); err != nil {
		t.Fatalf("iter.Close error: %v", err)
	}
}

func TestQueryPtr(t *testing.T) {
	db := mem.WithData(map[data.Kind][]data.Record{
		TestRecordKind: []data.Record{
			&TestRecord{
				Id:    "1",
				Name:  "third",
				Count: 3,
			},
			&TestRecord{
				Id:    "2",
				Name:  "second",
				Count: 2,
			},
			&TestRecord{
				Id:    "3",
				Name:  "first",
				Count: 1,
				True:  true,
				Ptr: &s{
					Foo: "yes",
				},
			},
		},
	})

	iter, err := db.Query(TestRecordKind).Order("Count").Execute()
	if err != nil {
		t.Fatalf("db.Query error: %v", err)
	}

	record := new(TestRecord)

	if got, want := iter.Next(record), true; got != want {
		t.Fatalf("iter.Next: got %t, want %t", got, want)
	}

	t.Logf("Loaded first record: %+v", record)

	if got, want := record.Name, "first"; got != want {
		t.Errorf("record.Name: got %q, want %q", got, want)
	}

	if got, want := record.True, true; got != want {
		t.Errorf("record.True: got %t, want %t", got, want)
	}

	if got, want := record.Ptr.Foo, "yes"; got != want {
		t.Errorf("record.Ptr.Foo: got %s, want %s", got, want)
	}

	if got, want := iter.Next(record), true; got != want {
		t.Fatalf("iter.Next: got %t, want %t", got, want)
	}

	t.Logf("Loaded second record: %+v", record)

	if got, want := record.Name, "second"; got != want {
		t.Errorf("record.Name: got %q, want %q", got, want)
	}

	if got, want := record.True, false; got != want {
		t.Fatalf("iter.Next: got %t, want %t", got, want)
	}

	if got, want := record.Ptr, (*s)(nil); got != want {
		t.Errorf("record.Ptr: got %v, want %v", got, want)
	}

	if got, want := iter.Next(record), true; got != want {
		t.Fatalf("iter.Next: got %t, want %t", got, want)
	}

	t.Logf("Loaded third record: %+v", record)

	if got, want := record.Name, "third"; got != want {
		t.Errorf("record.Name: got %q, want %q", got, want)
	}

	if got, want := record.Ptr, (*s)(nil); got != want {
		t.Errorf("record.Ptr: got %v, want %v", got, want)
	}

	if got, want := iter.Next(record), false; got != want {
		t.Fatalf("iter.Next: got %t, want %t", got, want)
	}

	if err := iter.Close(); err != nil {
		t.Fatalf("iter.Close error: %v", err)
	}
}
