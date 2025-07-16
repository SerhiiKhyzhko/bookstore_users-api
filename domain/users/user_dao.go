package users

import (
	"fmt"

	"github.com/SerhiiKhyzhko/bookstore_users-api/utils/errors"
)

//db mock
var usersDB = make(map[int64]*User)

func (user *User) Get() *errors.RestErr {
	result := usersDB[user.Id]
	if result == nil{
		return errors.NewNotFoundError(fmt.Sprintf("user %d not found", user.Id))
	}

	user.Id = result.Id
	user.FirstName = result.FirstName
	user.LastName = result.LastName
	user.Email = result.Email
	user.DateCreating = result.DateCreating

	return nil
}

func (user *User) Save() *errors.RestErr {
	if usersDB[user.Id] != nil {
		if usersDB[user.Id].Email == user.Email{
			return errors.NewBadRequestError(fmt.Sprintf("email %d already registered", user.Email))
		}
		return errors.NewBadRequestError(fmt.Sprintf("user %d already exists", user.Id))
	}

	usersDB[user.Id] = user

	return nil
}