package errs

import "errors"

// Module-level error
var (
	ErrInvalidArgument = errors.New("invalid argument")
)
