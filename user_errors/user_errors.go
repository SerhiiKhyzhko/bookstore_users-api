package user_errors

import "errors"

var (
	BadRequestErr     = errors.New("bad request")
	NotFoundErr       = errors.New("user not found")
	RequestTimeoutErr = errors.New("request timeout")
)
