package calendar

import "errors"

var (
	ErrUserNotFound  = errors.New("user not found")
	ErrEventNotFound = errors.New("event not found")
	ErrInvalidPeriod = errors.New("invalid period type")
)
