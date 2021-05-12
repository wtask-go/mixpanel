package errs

import "errors"

// Module-level error
var (
	ErrResponseInvalidData    = errors.New("invalid data")
	ErrResponseInvalidContent = errors.New("invalid content")
	ErrRequestFailed          = errors.New("request failed")
	ErrUnknown                = errors.New("unknown error")
)
