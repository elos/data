package data

import (
	"errors"
)

var NotFoundError = errors.New("Database error: record not found")
var InvalidIDError = errors.New("Database error: invalid id")
