package user

import "errors"

var ErrUserNotFound = errors.New("user not found")
var ErrUserExists = errors.New("user is already exists")
