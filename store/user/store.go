package user

import (
	"errors"

	"github.com/SemmiDev/todo-app/model"
)

var (
	ErrUserAlreadyExists = errors.New("user already exists")
	ErrUserNotFound      = errors.New("user not found")
	ErrUsersIsEmpty      = errors.New("users is empty")
)

type UserStore interface {
	Save(user *model.User) error
	Get(username string) (*model.User, error)
	ExistsByEmail(email string) bool
	ExistsByUsername(username string) bool
}
