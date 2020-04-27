package pkg

import "errors"

var (
	// ErrInvalidProviderResponse occurs when a given provider could not product a
	// valid response due to internal error like "InvalidSQL", "Invalid Transaction", etc.
	ErrInvalidProviderResponse = errors.New("invalid provider response")
	// ErrRecordNotFound the given query does not return any record.
	ErrRecordNotFound = errors.New("record not found")
)
