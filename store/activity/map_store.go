package activity

import (
	"context"
	"log"
	"sync"

	"github.com/SemmiDev/todo-app/model"
)

type MapStore struct {
	m          sync.RWMutex
	activities map[string]*model.Activity
}

func NewMapStore() *MapStore {
	return &MapStore{
		m:          sync.RWMutex{},
		activities: make(map[string]*model.Activity),
	}
}

func (s *MapStore) Save(todo *model.Activity) error {
	s.m.Lock()
	defer s.m.Unlock()

	if _, ok := s.activities[todo.Id]; !ok {
		s.activities[todo.Id] = Copy(todo)
		return nil
	}
	return ErrActivityAlreadyExists
}

func (s *MapStore) Get(id string) (*model.Activity, error) {
	s.m.RLock()
	defer s.m.RUnlock()

	activity, ok := s.activities[id]
	if !ok {
		return nil, ErrActivityNotFound
	}
	return Copy(activity), nil
}

func (s *MapStore) GetIdByDate(date string) ([]string, error) {
	s.m.RLock()
	defer s.m.RUnlock()

	if len(s.activities) == 0 {
		return nil, ErrActivityIsEmpty
	}
	var ids []string
	for i, v := range s.activities {
		// MM-DD-YYYY
		format := v.GetCreatedAt().AsTime().Format("01-02-2006")
		if format == date {
			ids = append(ids, i)
		}
	}
	return ids, nil
}

func (s *MapStore) List() ([]*model.Activity, error) {
	s.m.RLock()
	defer s.m.RUnlock()

	if len(s.activities) == 0 {
		return nil, ErrActivityIsEmpty
	}
	return CopyAll(s.activities), nil
}

func (s *MapStore) Delete(id string) error {
	s.m.Lock()
	defer s.m.Unlock()

	if _, ok := s.activities[id]; !ok {
		return ErrActivityNotFound
	}
	delete(s.activities, id)
	return nil
}

func (s *MapStore) Update(todo *model.Activity) error {
	s.m.Lock()
	defer s.m.Unlock()

	if _, ok := s.activities[todo.Id]; !ok {
		return ErrActivityNotFound
	}

	todo.CreatedAt = s.activities[todo.Id].GetCreatedAt()
	s.activities[todo.Id] = Copy(todo)
	return nil
}

func (s *MapStore) Search(c context.Context, filter *model.SearchActivityFilter, found func(todo *model.Activity) error) error {
	s.m.RLock()
	defer s.m.RUnlock()

	for _, todo := range s.activities {
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

func isQualified(filter *model.SearchActivityFilter, todo *model.Activity) bool {
	return todo.GetDay() == filter.GetDay()
}
