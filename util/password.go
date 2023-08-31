package util

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

func HashPassWord(password string) (string, error) {
	hashPassWord, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("fail to hash password: %s", err)
	}
	return string(hashPassWord), err
}

func CheckPassWord(password string, hashedPassword string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}
