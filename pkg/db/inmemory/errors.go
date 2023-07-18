package inmemory

import "errors"

var NotFoundError = errors.New("not found")

var AlreadyExistsError = errors.New("already exists")

var MissingRequiredFieldsError = errors.New("missing required fields")
