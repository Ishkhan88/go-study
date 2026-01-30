package service

import "errors"

var (
	ErrBadInput = errors.New("bad input")
	ErrNotFound = errors.New("not found")
)
