package cryptoutils

import (
	"github.com/SerhiiKhyzhko/bookstore_users-api/utils/errors"
	"golang.org/x/crypto/bcrypt"
)

func GetBcrypt(input string) (string, *rest_errors.RestErr) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input), 12)
	if err != nil {
		return "", rest_errors.NewBadRequestError(err.Error())
	}
	return string(hashedPassword), nil
}

func ComparePassword(hashedPassword string, password string) *errors.RestErr {
	if err:= bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password)); err != nil{
		return rest_errors.NewBadRequestError(err.Error())
	}
	return nil
}