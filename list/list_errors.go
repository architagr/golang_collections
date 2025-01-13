package list

import "errors"

var (
	ErrInvalidIndex      = errors.New("invalid index")
	ErrDataNotFoundError = errors.New("data not found")
	ErrDoesNotImplement  = errors.New("object does not implement")
)
