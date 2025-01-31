package domain

import "errors"

var (
	ErrFetchingUser     = errors.New("ERR_FETCHING_USER")
	ErrHashingPassword  = errors.New("ERR_HASHING_PASSWORD")
	ErrUsernameExists   = errors.New("ERR_USERNAME_EXISTS")
	ErrCreatingUser     = errors.New("ERR_CREATING_USER")
	ErrUsernameNotFound = errors.New("ERR_USERNAME_NOT_FOUND")
	ErrTokenGeneration  = errors.New("ERR_TOKEN_GENERATION")
	ErrInvalidCategory  = errors.New("ERR_INVALID_CATEGORY")
	ErrCodeExists       = errors.New("ERR_CODE_EXISTS")
	ErrFetchingItem     = errors.New("ERR_FETCHING_ITEM")
	ErrItemNotFound     = errors.New("ERR_ITEM_NOT_FOUND")
	ErrUpdatingItem     = errors.New("ERR_UPDATING_ITEM")
	ErrMissingFields    = errors.New("ERR_MISSING_FIELDS")
)
