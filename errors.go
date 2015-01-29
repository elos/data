package data

import (
	"errors"
)

var (
	// ErrNotFound is used if a Record is not retrievable.
	ErrNotFound = errors.New("database error: record not found")

	// ErrInvalidID is used if a type assertion on a Records's ID fails
	// the DB's implementation.
	ErrInvalidID = errors.New("database error: invalid id")

	// ErrInvalidDBType is used if a DB decides a Record's DBType is
	// incompatible with the DB's DBType.
	// This is generally only an issue for SQL variants
	ErrInvalidDBType = errors.New("database error: invalid type")

	ErrUndefinedKind      = errors.New("undefined kind")
	ErrUndefinedLink      = errors.New("undefined link")
	ErrUndefinedLinkKind  = errors.New("undefined link mind")
	ErrInvalidSchema      = errors.New("invalid schema")
	ErrIncompatibleModels = errors.New("incompatible models")
)
