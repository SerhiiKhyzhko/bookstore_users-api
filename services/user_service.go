package services

import (
	"github.com/SerhiiKhyzhko/bookstore_users-api/domain/users"
	"github.com/SerhiiKhyzhko/bookstore_users-api/utils/errors"
)

func CreateUser(user users.User) (*users.User, *errors.RestErr) {
	return &user, nil
}