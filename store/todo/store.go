package todo

import (
	"context"
	"errors"

	"github.com/SemmiDev/todo-app/model"
)

var (
	ErrTodoAlreadyExists = errors.New("todo already exists")
	ErrTodoNotFound      = errors.New("todo not found")
	ErrTodosIsEmpty      = errors.New("todo is empty")
)

type TodoStore interface {
	Save(todo *model.Todo) error
	Get(id string) (*model.Todo, error)
	List() ([]*model.Todo, error)
	ListByActivityId(id string) ([]*model.Todo, error)
	ListByActivityIds(id []string) ([]*model.Todo, error)
	Delete(id string) error
	Update(todo *model.Todo) error
	Search(c context.Context, filter *model.SearchTodoFilter, found func(todo *model.Todo) error) error
}

func Copy(todo *model.Todo) *model.Todo {
	return &model.Todo{
		Id:          todo.Id,
		Title:       todo.Title,
		ActivityId:  todo.ActivityId,
		Description: todo.Description,
		Priority:    todo.Priority,
		Status:      todo.Status,
		CreatedAt:   todo.CreatedAt,
		UpdatedAt:   todo.CreatedAt,
	}
}

func CopyAll(todos map[string]*model.Todo) (data []*model.Todo) {
	for _, v := range todos {
		data = append(data, Copy(v))
	}
	return
}
