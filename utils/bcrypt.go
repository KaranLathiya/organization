package utils

import (
	"golang.org/x/crypto/bcrypt"
)

func Bcrypt(value string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(value), 14)
	return string(bytes), err
}

func CompareHashAndPassword(hashPassword []byte, password []byte) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashPassword), []byte(password))
	return err == nil
}
