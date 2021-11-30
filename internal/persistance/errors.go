package persistance

import "errors"

var ErrUserNotFound error = errors.New("user not found")
var ErrInvalidCredentials error = errors.New("invalid credential")
var ErrInvalidEmail error = errors.New("invalid email")
var ErrCannotFindMail error = errors.New("can not find email")
var ErrMailAlreadyExist error = errors.New("email is already used")
