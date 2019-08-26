package shortener

import "errors"

var (
	// ErrRedirectNotFound error
	ErrRedirectNotFound = errors.New("Redirect Not Found")
	// ErrRedirectInvalid error
	ErrRedirectInvalid = errors.New("Redirect Invalid")
)
