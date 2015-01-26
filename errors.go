package data

import (
	"errors"
)

var (
	ErrNotFound    = errors.New("database error: record not found")
	ErrInvalidID   = errors.New("database error: invalid id")
	ErrInvalidtype = errors.New("database error: invalid type")
)
