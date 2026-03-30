package api

import "errors"

var (
	ErrBadPayload     = errors.New("bad payload")
	ErrBadCredentials = errors.New("bad credentials")
	ErrInternal       = errors.New("internal error")
	ErrBadPath        = errors.New("not provided proper path values")
	ErrUnauthorized   = errors.New("not authorized")
)
