package transfer

import (
	"encoding/json"

	"github.com/elos/data"
)

// Returns a map of form:
//		{ <db.Kind>: <db.Model>}
func Map(r data.Record) data.KindMap {
	return data.KindMap{
		r.Kind(): r,
	}
}

// Transfers json based struct fields from this to that
func TransferAttrs(this interface{}, that interface{}) error {
	bytes, err := json.Marshal(this)
	if err != nil {
		return err
	}
	return json.Unmarshal(bytes, that)
}

// Unmarshal json version of attr map into a record
func Unmarshal(attrs data.AttrMap, r data.Record) (data.Record, error) {
	bytes, err := json.Marshal(attrs)
	if err != nil {
		return r, err
	}
	return r, json.Unmarshal(bytes, r)
}
