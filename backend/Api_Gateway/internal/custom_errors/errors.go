package errors

import "errors"

var (
	ErrUserExists      = errors.New("user already exists")
	ErrUserNotFound    = errors.New("user not found")
	ErrAppNotFound     = errors.New("app not found")
	ErrChannelsEmpty   = errors.New("channels list cannot be empty")
	ErrChannelExists   = errors.New("channel already added for this user")
	ErrChannelNotFound = errors.New("channel not found")
)
