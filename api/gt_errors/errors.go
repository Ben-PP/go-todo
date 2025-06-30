package gterrors

import "errors"

var ErrConfigLoadFailed = errors.New("failed to load config")
var ErrShouldNotHappen = errors.New("this should not happen")
var ErrJwtRefreshReuse = errors.New("refresh jwt reuse")
var ErrPasswordUnsatisfied = errors.New("password criteria not met")