package domain

import "errors"

var (
	ErrLinkNotFound   = errors.New("link not found")
	ErrInvalidURL     = errors.New("invalid URL")
	ErrInvalidShortID = errors.New("invalid short ID")
	ErrStorage        = errors.New("storage error")
	ErrDuplicate      = errors.New("duplicate link")
)
