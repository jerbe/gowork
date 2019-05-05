package gowork

import "errors"

var (
	// cap less than zero
	ErrCapLessThanZero = errors.New("cap can't not less than zero")

	ErrJobQueueNil = errors.New("job queue was nil")

	ErrWorkPoolNil = errors.New("work pool was nil")
)
