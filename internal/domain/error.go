package domain

import "errors"

var ErrCodeExists = errors.New("code already exists")
var ErrInvalidCategory = errors.New("invalid category")
var ErrItemNotFound = errors.New("item not found")
var ErrMissingFields = errors.New("missing required fields")
