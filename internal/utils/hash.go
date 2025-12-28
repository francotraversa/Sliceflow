package utils

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

func CheckPassword(hash, plain string) error {
	if bcrypt.CompareHashAndPassword([]byte(hash), []byte(plain)) != nil {
		return fmt.Errorf("Incorrect Password")
	}
	return nil
}
