package auth

import "errors"

var ErrUserNotFound = errors.New("User not found")
var ErrUserAlreadyExists = errors.New("User already exists")
var ErrWrongPassword = errors.New("Wrong password")