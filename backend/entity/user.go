package entity

import (
	"golang.org/x/crypto/bcrypt"
	"html"
	"strings"
)

const (
	RoleAdmin  = 0
	RoleUser   = 1
	RoleBarber = 2
)

type Account struct {
	ID          int    `json:"id,omitempty"`
	Username    string `json:"username"`
	Email       string `json:"email"`
	Password    string `json:"password,omitempty"`
	Role        int    `json:"role,omitempty"`
	PhoneNumber string `json:"phoneNumber"`
	FullName    string `json:"fullName"`
	About       string `json:"about"`
	Avatar      string `json:"avatar"`
	CreatedAt   int64  `json:"created_at,omitempty"`
	UpdatedAt   int64  `json:"updated_at,omitempty"`
}

type Role struct {
	ID   int    `json:"id,omitempty"`
	Name string `json:"name"`
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
