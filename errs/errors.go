package errs

import "errors"

// Module-level error
var (
	ErrResponseInvalidContent = errors.New("invalid content")
	ErrRequestFailed          = errors.New("request failed")
	ErrUnknown                = errors.New("unknown error")
)
