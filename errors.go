package mongo

import "errors"

var (
	// ErrInvalidConfig indicates validation failed on Config.
	ErrInvalidConfig = errors.New("mongo: invalid config")
)
