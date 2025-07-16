package services

import (
	"github.com/SerhiiKhyzhko/bookstore_users-api/domain/users"
	"github.com/SerhiiKhyzhko/bookstore_users-api/utils/errors"
)

func CreateUser(user users.User) (*users.User, *errors.RestErr) {
	if err := user.ValidateEmail(); err != nil{
		return nil, err
	}

	if err := user.Save(); err != nil {
		return nil, err
	}

	return &user, nil
}

func GetUser(userId int64) (*users.User, *errors.RestErr) {
	result := users.User{Id: userId}
	if err := result.Get(); err != nil {
		return nil, err
	}
	return &result, nil
}