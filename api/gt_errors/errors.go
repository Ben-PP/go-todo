package gterrors

import "errors"

var ErrConfigLoadFailed = errors.New("failed to load config")
var ErrForbidden = errors.New("forbidden")
var ErrJwtRefreshReuse = errors.New("refresh jwt reuse")
var ErrNotFound = errors.New("resource not found")
var ErrPasswordUnsatisfied = errors.New("password criteria not met")
var ErrShouldNotHappen = errors.New("this should not happen")
var ErrUniqueViolation = errors.New("allready exists")
var ErrUsernameUnsatisfied = errors.New("username criteria not met")