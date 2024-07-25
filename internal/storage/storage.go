package storage

import "errors"

var (
	ErrorUserExists   = errors.New("user already exists")
	ErrorUserNotFound = errors.New("user not found")
	ErrorAppNotFound  = errors.New("app not found")
)
