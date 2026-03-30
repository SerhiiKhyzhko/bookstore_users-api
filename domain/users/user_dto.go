package users

import (
	"strings"

	"github.com/SerhiiKhyzhko/bookstore_users-api/user_errors"
)

const(
	StatusActive = "active"
)

type User struct {
	Id           int64  `json:"id"`
	FirstName    string `json:"first_name"`
	LastName     string `json:"last_name"`
	Email        string `json:"email"`
	DateCreating string `json:"date_creating"`
	Status 		 string `json:"status"`
	Password 	 string `json:"password"`
}

type Users []User

func (user *User) Validate() error {
	user.FirstName = strings.TrimSpace(user.FirstName)
	user.LastName = strings.TrimSpace(user.LastName)
	
	user.Email = strings.TrimSpace(strings.ToLower(user.Email))
	if user.Email == ""{
		return fmt.Errorf("%w: invalid email", user_errors.BadRequestErr)
	}
	user.Password = strings.TrimSpace(user.Password)
	if len(user.Password) < 4 {
		return fmt.Errorf("%w: invalid password. Password has to be at least 4 symbols", user_errors.BadRequestErr)
	}
	return nil
}

type PartialUser struct {
	Id           int64   `json:"id"`
	FirstName    *string `json:"first_name"`
	LastName     *string `json:"last_name"`
	Email        *string `json:"email"`
	Status 		 *string `json:"status"`
}