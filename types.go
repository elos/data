package data

import "encoding/json"

// DBType is defined by each implementation of a DB,
// so that the DB can be semi-identifiable and implement
// compatibility checking against a Record's declared DBType
type DBType string

// Kind is a record's table name of collection name.
// It should correspond to the model's name, generally lowercase.
type Kind string

// A KindMap represents a mapping from
// Record Kind to Record object, useful
// for protocols which follow:
//	{ <kind>: { ... info ... } }
type KindMap map[Kind]Record

/*
	Returns a map like:
	{ user: { Name: "Nick Landolfi"} }
	of form:
	{ <db.Kind>: <db.Model>}
*/
func Map(r Record) KindMap {
	return KindMap{
		r.Kind(): r,
	}
}

// An AttrMap is the type used to
// populate a Record's fields.
type AttrMap map[string]interface{}

func TransferAttrs(this interface{}, that interface{}) error {
	bytes, err := json.Marshal(this)
	if err != nil {
		return nil
	}
	return json.Unmarshal(bytes, that)
}
