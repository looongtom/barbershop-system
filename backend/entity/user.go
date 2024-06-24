package entity

import (
	"golang.org/x/crypto/bcrypt"
	"html"
	"strings"
)

type User struct {
	ID       int    `json:"id,omitempty"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password,omitempty"`
	Roles    string `json:"roles,omitempty"`
}

func Hash(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}
func CheckPasswordHash(hashedPassword, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}
func Santize(data string) string {
	data = html.EscapeString(strings.TrimSpace(data))
	return data
}
