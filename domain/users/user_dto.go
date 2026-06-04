package users

import (
	"fmt"
	"strings"
	"unicode/utf8"

	"github.com/SerhiiKhyzhko/bookstore_users-api/v2/user_errors"
)

const (
	StatusActive = "active"
)

type User struct {
	Id           int64  `json:"id"`
	FirstName    string `json:"first_name"`
	LastName     string `json:"last_name"`
	Email        string `json:"email"`
	DateCreating string `json:"date_creating"`
	Status       string `json:"status"`
	Password     string `json:"password"`
}

type Users []User

func (user *User) Validate() error {
	user.FirstName = strings.TrimSpace(user.FirstName)
	user.LastName = strings.TrimSpace(user.LastName)
	user.Email = strings.TrimSpace(strings.ToLower(user.Email))

	if count := utf8.RuneCountInString(user.FirstName); count < 1 {
		return fmt.Errorf("%w: first name is required", user_errors.BadRequestErr)
	} else if count > 45 {
		return fmt.Errorf("%w: first name is too long (max 45)", user_errors.BadRequestErr)
	}

	if count := utf8.RuneCountInString(user.LastName); count < 1 {
		return fmt.Errorf("%w: last name is required", user_errors.BadRequestErr)
	} else if count > 45 {
		return fmt.Errorf("%w: last name is too long (max 45)", user_errors.BadRequestErr)
	}

	if count := utf8.RuneCountInString(user.Email); count < 1 {
		return fmt.Errorf("%w: email is required", user_errors.BadRequestErr)
	} else if count > 45 {
		return fmt.Errorf("%w: email is too long (max 45)", user_errors.BadRequestErr)
	}

	if utf8.RuneCountInString(user.Password) < 4 {
		return fmt.Errorf("%w: password must be at least 4 characters", user_errors.BadRequestErr)
	}

	return nil
}

type PartialUser struct {
	Id        int64   `json:"id"`
	FirstName *string `json:"first_name"`
	LastName  *string `json:"last_name"`
	Email     *string `json:"email"`
	Status    *string `json:"status"`
}
