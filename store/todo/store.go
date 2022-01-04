package todo

import (
	"context"
	"errors"
	"github.com/SemmiDev/todo-app/proto"
)

var (
	ErrTodoAlreadyExists = errors.New("todo already exists")
	ErrTodoNotFound      = errors.New("todo not found")
	ErrTodosIsEmpty      = errors.New("todo is empty")
)

type TodoStore interface {
	Save(todo *proto.Todo) error
	Get(id string) (*proto.Todo, error)
	List() ([]*proto.Todo, error)
	ListByActivityId(id string) ([]*proto.Todo, error)
	ListByActivityIds(id []string) ([]*proto.Todo, error)
	Delete(id string) error
	Update(todo *proto.Todo) error
	Search(c context.Context, filter *proto.SearchTodoFilter, found func(todo *proto.Todo) error) error
}

func Copy(todo *proto.Todo) *proto.Todo {
	return &proto.Todo{
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

func CopyAll(todos map[string]*proto.Todo) (data []*proto.Todo) {
	for _, v := range todos {
		data = append(data, Copy(v))
	}
	return
}
