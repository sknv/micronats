package models

import (
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID           string
	Login        string
	PasswordHash string
}

func (u *User) Authenticate(password string) error {
	return bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(password))
}
