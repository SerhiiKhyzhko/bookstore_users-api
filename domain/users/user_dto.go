package users

import (
	"fmt"
	"strings"

	"github.com/SerhiiKhyzhko/bookstore_users-api/utils/errors"
)

const(
	StatusActive = "active"
)

type User struct {
	Id           int64    `json:"id"`
	FirstName    string `json:"first_name"`
	LastName     string `json:"last_name"`
	Email        string `json:"email"`
	DateCreating string `json:"date_creating"`
	Status 		 string `json:"status"`
	Password 	 string `json:"-"`
}

type Users []User

func (user *User) Validate() *errors.RestErr {
	user.FirstName = strings.TrimSpace(user.FirstName)
	user.LastName = strings.TrimSpace(user.LastName)
	
	user.Email = strings.TrimSpace(strings.ToLower(user.Email))
	if user.Email == ""{
		return errors.NewBadRequestError("invalid email")
	}
	user.Password = strings.TrimSpace(user.Password)
	fmt.Println(user.Password, len(user.Password))
	if len(user.Password) < 4 {
		return errors.NewBadRequestError("invalid password. Password has to be at least 4 symbols")
	}
	return nil
}