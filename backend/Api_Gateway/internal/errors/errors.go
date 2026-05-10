package errs

import "errors"

var (
	ErrUserExists      = errors.New("user already exists")
	ErrIsEmpty         = errors.New("data is empty")
	ErrUserNotFound    = errors.New("user not found")
	ErrAppNotFound     = errors.New("app not found")
	ErrCacheMiss       = errors.New("cache miss")
	ErrChannelsEmpty   = errors.New("channels list cannot be empty")
	ErrChannelExists   = errors.New("channel already added for this user")
	ErrChannelNotFound = errors.New("channel not found")
)
