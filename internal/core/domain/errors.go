package domain

import "errors"

var ErrNotFound = errors.New("code not found")
var ErrCodeExists = errors.New("code already exists")
