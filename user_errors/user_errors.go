package user_errors

import "errors"

var (
	NotFoundErr       = errors.New("user not found")
	RequestTimeoutErr = errors.New("request timeout")
	BadRequestErr     = errors.New("bad request")
)
