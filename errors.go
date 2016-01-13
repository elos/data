package data

import "fmt"

// formatError prepends a package identifier to all error messages
func formatError(s string) error {
	return fmt.Errorf("data Error: %s", s)
}

var (
	// ErrNotFound indicates that a record which was searched
	// or queried for could not be retrieved.
	//
	// Use ErrNotFound to indicate the request could not be fulfilled
	// in the case of PopulateByID and PopulateByField. For an empty query,
	// rather, simply return an empty iterator.
	ErrNotFound = formatError("record not found")

	// ErrNoConnection indicates some sort of connection error.
	//
	// Use ErrNoConnection for DB implementations which traverse the network,
	// and in the case that they fail. It should be used for all network
	// failures.
	ErrNoConnection = formatError("no database connection")

	// ErrInvalidID indicates that a record given to a DB has an invalid ID
	//
	// We need this error of the rather loose typing on the types of ID
	// which the DB interface supports. Some databases use numerical ids whereas
	// others use UUIDs.
	//
	// Use ErrInvalidID to indicate that the ID of a record handed to a DB
	// had an ID that could not be satisfactorily parsed.
	ErrInvalidID = formatError("invalid id")

	// ErrAccessDenial indicates that the DB does that have access to a particular
	// table, or record, or query mechanism.
	//
	// Use ErrAccessDenial to indicate that some sort of AccessDenial was encountered.
	// You may find it effective to wrap a data.DB in some sort of structure which
	// maintains the access rules for you schema. In this case, said structure could
	// use ErrAccessDenial to reject access to a client.
	ErrAccessDenial = formatError("access denied")
)
