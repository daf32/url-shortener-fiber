package core_server

import "errors"

var ErrNotFound = errors.New("code not found")
var ErrCodeExists = errors.New("code already exists")
var ErrInvalidArgument = errors.New("invalid argument")
