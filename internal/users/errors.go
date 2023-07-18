package users

import "errors"

var UserNotFoundError = errors.New("user not found")

var UserAlreadyExistsError = errors.New("user with this username already exists")

var UnknownError = errors.New("unknown error")
