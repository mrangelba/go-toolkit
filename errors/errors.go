package errors

import "errors"

var ErrRecordNotFound = errors.New("record not found")
var ErrInvalidRequest = errors.New("invalid request")
var ErrInvalidID = errors.New("invalid id")
