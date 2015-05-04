package data

import (
	"errors"
	"fmt"
)

func formatError(s string) error {
	return errors.New(fmt.Sprintf("data error: %s", s))
}

var (
	// ErrNotFound is when PopulateByID or PopulateByField
	// could not be fulfilled
	ErrNotFound = formatError("record not found")

	// ErrNoConnection is used if a DB drive has not yet
	// connected to, or has lost connection with the DBMS
	ErrNoConnection = formatError("no database connection")

	// ErrInvalidID is used if a Record given to a driver
	// does not have a validly parseable id
	ErrInvalidID = formatError("invalid id")

	// ErrAccessDenial is used if a driver encounters some sort
	// of access restriction from the DBMS. It is also used if
	// an implementation of data.Access rejects access to a Client
	ErrAccessDenial = formatError("access denied")
)
