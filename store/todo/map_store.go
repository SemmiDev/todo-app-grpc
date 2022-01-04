package todo

import (
	"context"
	"github.com/SemmiDev/todo-app/proto"
	"log"
	"sync"
)

type MapStore struct {
	m     sync.RWMutex
	todos map[string]*proto.Todo
}

func NewMapStore() *MapStore {
	return &MapStore{
		m:     sync.RWMutex{},
		todos: make(map[string]*proto.Todo),
	}
}

func (s *MapStore) Save(todo *proto.Todo) error {
	s.m.Lock()
	defer s.m.Unlock()

	if _, ok := s.todos[todo.Id]; !ok {
		s.todos[todo.Id] = Copy(todo)
		return nil
	}
	return ErrTodoAlreadyExists
}

func (s *MapStore) Get(id string) (*proto.Todo, error) {
	s.m.RLock()
	defer s.m.RUnlock()

	todo, ok := s.todos[id]
	if !ok {
		return nil, ErrTodoNotFound
	}
	return Copy(todo), nil
}

func (s *MapStore) List() ([]*proto.Todo, error) {
	s.m.RLock()
	defer s.m.RUnlock()

	if len(s.todos) == 0 {
		return nil, ErrTodosIsEmpty
	}
	return CopyAll(s.todos), nil
}

func (s *MapStore) ListByActivityId(id string) ([]*proto.Todo, error) {
	s.m.RLock()
	defer s.m.RUnlock()

	if len(s.todos) == 0 {
		return nil, ErrTodosIsEmpty
	}

	var todos []*proto.Todo
	for _, v := range s.todos {
		if v.GetActivityId() == id {
			todos = append(todos, v)
		}
	}
	return todos, nil
}

func (s *MapStore) ListByActivityIds(ids []string) ([]*proto.Todo, error) {
	s.m.RLock()
	defer s.m.RUnlock()

	if len(s.todos) == 0 {
		return nil, ErrTodosIsEmpty
	}

	var todos []*proto.Todo
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
	filter *proto.SearchTodoFilter,
	found func(todo *proto.Todo) error,
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

func isQualified(filter *proto.SearchTodoFilter, todo *proto.Todo) bool {
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

func (s *MapStore) Update(todo *proto.Todo) error {
	s.m.Lock()
	defer s.m.Unlock()

	if _, ok := s.todos[todo.Id]; !ok {
		return ErrTodoNotFound
	}

	todo.CreatedAt = s.todos[todo.Id].GetCreatedAt()
	s.todos[todo.Id] = Copy(todo)
	return nil
}
