package transfer_test

import (
	"testing"

	"github.com/elos/data/transfer"
)

type s struct {
	Foo string
}

type T struct {
	One string
	Two int

	Other *s
}

func TestTransferAttrsToMap(t *testing.T) {
	x := &T{
		One: "one",
		Two: 2,
		Other: &s{
			Foo: "yes",
		},
	}

	m := make(map[string]interface{})

	transfer.TransferAttrs(x, &m)

	one, ok := m["One"]
	if got, want := ok, true; got != want {
		t.Fatalf("_, ok := m[\"One\"]: got %t, want %t", got, want)
	}
	if got, want := one.(string), "one"; got != want {
		t.Fatalf("one.(string): got %q, want %q", got, want)
	}

	two, ok := m["Two"]
	if got, want := ok, true; got != want {
		t.Fatalf("_, ok := m[\"Two\"]: got %t, want %t", got, want)
	}
	if got, want := two.(float64), float64(2); got != want {
		t.Fatalf("two.(string): got %f, want %f", got, want)
	}

	other, ok := m["Other"]
	if got, want := ok, true; got != want {
		t.Fatalf("_, ok := m[\"Other\"]: got %t, want %t", got, want)
	}
	foo, ok := other.(map[string]interface{})["Foo"]
	if got, want := ok, true; got != want {
		t.Fatalf("_, ok := m[\"Other\"][\"Foo\"]: got %t, want %t", got, want)
	}
	if got, want := foo.(string), "yes"; got != want {
		t.Fatalf("foo.(string): got %q, want %q", got, want)
	}
}

func TestTransferAttrsToStruct(t *testing.T) {
	from := &T{
		One: "one",
		Two: 2,
		Other: &s{
			Foo: "yes",
		},
	}

	to := &T{}

	transfer.TransferAttrs(from, to)

	if got, want := to.One, "one"; got != want {
		t.Errorf("to.One: got %q, want %q", got, want)
	}

	if got, want := to.Two, 2; got != want {
		t.Errorf("to.Two: got %f, want %f", got, want)
	}

	if got, want := to.Other.Foo, "yes"; got != want {
		t.Errorf("to.Other.Foo: got %q, want %q", got, want)
	}
}
