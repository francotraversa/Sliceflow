package utils

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

func CheckPassword(hash, plain string) error {
	if bcrypt.CompareHashAndPassword([]byte(hash), []byte(plain)) != nil {
		return fmt.Errorf("Something went wrong, please try again later")
	}
	return nil
}
