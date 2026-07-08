package biz

import "errors"

// Sentinel errors; the data layer translates driver errors into these so
// no layer above it imports the ORM.
var (
	ErrNotFound = errors.New("record not found")
)
