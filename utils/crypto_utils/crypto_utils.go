package cryptoutils

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

func GetBcrypt(input string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input), 12)
	if err != nil {
		return "", fmt.Errorf("error when trying to convert value: %w", err)
	}
	return string(hashedPassword), nil
}

func ComparePassword(hashedPassword string, password string) error {
	if err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password)); err != nil {
		return fmt.Errorf("Process failed: %w", err)
	}
	return nil
}
