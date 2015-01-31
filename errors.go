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

	ErrUndefinedKind      = errors.New("data error: undefined kind")
	ErrUndefinedLink      = errors.New("data error: undefined link")
	ErrUndefinedLinkKind  = errors.New("data error: undefined link mind")
	ErrInvalidSchema      = errors.New("data error: invalid schema")
	ErrIncompatibleModels = errors.New("data error: incompatible models")
)

type AttrError struct {
	AttrName string
	What     string
}

func (e AttrError) Error() string {
	return fmt.Sprintf("attribute %v must %v", e.AttrName, e.What)
}

func NewAttrError(a string, w string) AttrError {
	return AttrError{
		AttrName: a,
		What:     w,
	}
}
