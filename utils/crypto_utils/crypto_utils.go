package cryptoutils

import (
	"github.com/SerhiiKhyzhko/bookstore_users-api/utils/errors"
	"golang.org/x/crypto/bcrypt"
)

func GetBcrypt(input string) (string, *errors.RestErr) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input), 12)
	if err != nil {
		return "", errors.NewBadRequestError(err.Error())
	}
	return string(hashedPassword), nil
}