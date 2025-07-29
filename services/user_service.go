package services

import (
	"github.com/SerhiiKhyzhko/bookstore_users-api/domain/users"
	dateutils "github.com/SerhiiKhyzhko/bookstore_users-api/utils/date_utils"
	"github.com/SerhiiKhyzhko/bookstore_users-api/utils/errors"
)

func CreateUser(user users.User) (*users.User, *errors.RestErr) {
	if err := user.ValidateEmail(); err != nil{
		return nil, err
	}

	user.Status = users.StatusActive
	user.DateCreating = dateutils.GetNowDbFormat()
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

func UpdateUser(user users.User) (*users.User, *errors.RestErr) {
	current, err := GetUser(user.Id) 
	if err != nil {
		return nil, err
	}

	current.FirstName = user.FirstName
	current.LastName = user.LastName
	current.Email = user.Email
	

	if err = current.Update(); err != nil {
		return nil, err
	}

	return current, nil
}

func PartialUpdateUser(user users.User) (*users.User, *errors.RestErr) {
	current, err := GetUser(user.Id) 
	if err != nil {
		return nil, err
	}

	if user.FirstName != "" {
		current.FirstName = user.FirstName
	}
	if user.LastName != "" {
		current.LastName = user.LastName
	}
	if user.Email != "" {
		current.Email = user.Email
	}

	if err = current.Update(); err != nil {
		return nil, err
	}

	return current, nil
}

func DeleteUser(userId int64) *errors.RestErr {
	user := &users.User{Id: userId}
	return user.Delete()
}

func Search(status string) ([]users.User, *errors.RestErr) {
	dao := &users.User{}
	users, err := dao.FindByStatus(status)

	if err != nil {
		return nil, err
	}
	return  users, nil
}