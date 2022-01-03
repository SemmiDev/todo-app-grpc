package todo

import (
	"context"
	"log"
	"sync"

	"github.com/SemmiDev/todo-app/model"
)

type MapStore struct {
	m     sync.RWMutex
	todos map[string]*model.Todo
}

func NewMapStore() *MapStore {
	return &MapStore{
		m:     sync.RWMutex{},
		todos: make(map[string]*model.Todo),
	}
}

func (s *MapStore) Save(todo *model.Todo) error {
	s.m.Lock()
	defer s.m.Unlock()

	if _, ok := s.todos[todo.Id]; !ok {
		s.todos[todo.Id] = Copy(todo)
		return nil
	}
	return ErrTodoAlreadyExists
}

func (s *MapStore) Get(id string) (*model.Todo, error) {
	s.m.RLock()
	defer s.m.RUnlock()

	todo, ok := s.todos[id]
	if !ok {
		return nil, ErrTodoNotFound
	}
	return Copy(todo), nil
}

func (s *MapStore) List() ([]*model.Todo, error) {
	s.m.RLock()
	defer s.m.RUnlock()

	if len(s.todos) == 0 {
		return nil, ErrTodosIsEmpty
	}
	return CopyAll(s.todos), nil
}

func (s *MapStore) ListByActivityId(id string) ([]*model.Todo, error) {
	s.m.RLock()
	defer s.m.RUnlock()

	if len(s.todos) == 0 {
		return nil, ErrTodosIsEmpty
	}

	var todos []*model.Todo
	for _, v := range s.todos {
		if v.GetActivityId() == id {
			todos = append(todos, v)
		}
	}
	return todos, nil
}

func (s *MapStore) ListByActivityIds(ids []string) ([]*model.Todo, error) {
	s.m.RLock()
	defer s.m.RUnlock()

	if len(s.todos) == 0 {
		return nil, ErrTodosIsEmpty
	}

	var todos []*model.Todo
	for _, id := range ids {
		for _, v := range s.todos {
			if v.GetActivityId() == id {
				todos = append(todos, v)
			}
		}
	}
	return todos, nil
}

func (s *MapStore) Search(
	c context.Context,
	filter *model.SearchTodoFilter,
	found func(todo *model.Todo) error,
) error {
	s.m.RLock()
	defer s.m.RUnlock()

	for _, todo := range s.todos {
		if c.Err() == context.Canceled || c.Err() == context.DeadlineExceeded {
			log.Print("context is cancelled")
			return nil
		}

		if isQualified(filter, todo) {
			err := found(Copy(todo))
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func isQualified(filter *model.SearchTodoFilter, todo *model.Todo) bool {
	if todo.GetPriority() != filter.GetPriority() {
		return false
	}
	if todo.GetStatus() != filter.GetStatus() {
		return false
	}
	return true
}

func (s *MapStore) Delete(id string) error {
	s.m.Lock()
	defer s.m.Unlock()

	if _, ok := s.todos[id]; !ok {
		return ErrTodoNotFound
	}
	delete(s.todos, id)
	return nil
}

func (s *MapStore) Update(todo *model.Todo) error {
	s.m.Lock()
	defer s.m.Unlock()

	if _, ok := s.todos[todo.Id]; !ok {
		return ErrTodoNotFound
	}

	todo.CreatedAt = s.todos[todo.Id].GetCreatedAt()
	s.todos[todo.Id] = Copy(todo)
	return nil
}
