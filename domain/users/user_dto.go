package users

import (
	"strings"

	"github.com/SerhiiKhyzhko/bookstore_users-api/utils/errors"
)

type User struct {
	Id           int64    `json:"id"`
	FirstName    string `json:"first_name"`
	LastName     string `json:"last_name"`
	Email        string `json:"email"`
	DateCreating string `json:"date_creating"`
}

func (user *User) ValidateEmail() *errors.RestErr {
	user.Email = strings.TrimSpace(strings.ToLower(user.Email))
	if user.Email == ""{
		return errors.NewBadRequestError("invalid email")
	}
	return nil
}