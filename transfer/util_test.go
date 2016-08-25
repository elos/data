package transfer_test

import (
	"testing"

	"github.com/elos/data/transfer"
)

type T struct {
	One string
	Two int
}

func TestTransferAttrs(t *testing.T) {
	x := &T{
		One: "one",
		Two: 2,
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
}
