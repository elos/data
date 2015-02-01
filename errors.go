package data

import (
	"errors"
	"fmt"
)

var (
	// ErrNotFound is used if a Record is not retrievable.
	ErrNotFound = errors.New("data error: record not found")

	// ErrInvalidID is used if a type assertion on a Records's ID fails
	// the DB's implementation.
	ErrInvalidID = errors.New("data error: invalid id")

	// ErrInvalidDBType is used if a DB decides a Record's DBType is
	// incompatible with the DB's DBType.
	// This is generally only an issue for SQL variants
	ErrInvalidDBType = errors.New("data error: invalid type")

	// ErrInvalidSchema is used if the schema can not be validated
	ErrInvalidSchema = errors.New("data error: invalid schema")

	// ErrUndefinedKind is used if a DB or Schema or Store can't
	// verify or recognize the kind supplied by a Record
	ErrUndefinedKind = errors.New("data error: undefined kind")

	// ErrUndefinedLink is used if a link cannot be made
	// as it is not defined in the schema
	ErrUndefinedLink = errors.New("data error: undefined link")

	// ErrUndefinedLinkKind is used if a LinkKind has not been defined
	ErrUndefinedLinkKind = errors.New("data error: undefined link mind")

	// ErrIncompatibleModels is used if two models don't have the same DBType
	ErrIncompatibleModels = errors.New("data error: incompatible models")
)

// An AttrError is used when a model fails validation,
// the model's Valid() function should return false, AttrError
// Use the NewAttrError function to create AttrErrors
type AttrError struct {
	AttrName string
	What     string
}

// Error() is defined so that go can print the error nicely
func (e AttrError) Error() string {
	return fmt.Sprintf("attribute %v must %v", e.AttrName, e.What)
}

// NewAttrError constructs a new instance of an attr error.
//	var err AttError = NewAttrError("name", "be present")
// The error will print: "attribute [first string] must [second string]"
func NewAttrError(a string, w string) AttrError {
	return AttrError{
		AttrName: a,
		What:     w,
	}
}

type LinkError struct {
	this Model
	that Model
	link Link
}

func (e *LinkError) Error() string {
	return fmt.Sprintf("%T could not be linked to %T according to %+v", e.this, e.that, e.link)
}

func NewLinkError(ts Model, tt Model, l Link) *LinkError {
	return &LinkError{
		this: ts,
		that: tt,
		link: l,
	}
}
