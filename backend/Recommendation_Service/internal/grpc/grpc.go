package grpc

import "errors"

var (
	ErrUserExists      = errors.New("user already exists")
	ErrUserIDEmpty     = errors.New("user id is empty")
	ErrUserNotFound    = errors.New("user not found")
	ErrAppNotFound     = errors.New("app not found")
	ErrChannelsEmpty   = errors.New("channels list cannot be empty")
	ErrChannelExists   = errors.New("channel already added for this user")
	ErrChannelNotFound = errors.New("channel not found")
)
