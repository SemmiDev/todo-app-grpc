package model

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

type User struct {
	Username       string `json:"username"`
	Email          string `json:"email"`
	HashedPassword string `json:"-"`
	Role           string `json:"-"`
}

func NewUser(username, email, password, role string) (*User, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("cannot hash password: %w", err)
	}
	user := User{
		Username:       username,
		Email:          email,
		HashedPassword: string(hashedPassword),
		Role:           role,
	}
	return &user, nil
}

func (user *User) IsCorrectPassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(user.HashedPassword), []byte(password))
	return err == nil
}

func (user *User) Clone() *User {
	return &User{
		Username:       user.Username,
		Email:          user.Email,
		HashedPassword: user.HashedPassword,
		Role:           user.Role,
	}
}
