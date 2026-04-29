package cryptoutils

import (
	"errors"

	"golang.org/x/crypto/bcrypt"
)

func GetBcrypt(input string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input), 12)
	if err != nil {
		return "", errors.New("error when trying to hash password")
	}
	return string(hashedPassword), nil
}

func ComparePassword(hashedPassword string, password string) error {
	if err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password)); err != nil {
		return errors.New("invalid credentials")
	}
	return nil
}
