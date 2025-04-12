package connection

import "errors"

var (
	// ErrUnauthorized is returned when authentication fails
	ErrUnauthorized = errors.New("unauthorized connection")

	// ErrInvalidToken is returned when the token is invalid
	ErrInvalidToken = errors.New("invalid token")

	// ErrInvalidSession is returned when the session is invalid
	ErrInvalidSession = errors.New("invalid session")
)
