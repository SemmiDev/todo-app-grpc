package user

import (
	"sync"

	"github.com/SemmiDev/todo-app/model"
)

type MapStore struct {
	mutex sync.RWMutex
	users map[string]*model.User
}

func NewMapStore() *MapStore {
	return &MapStore{
		users: make(map[string]*model.User),
	}
}

func (store *MapStore) Save(user *model.User) error {
	store.mutex.Lock()
	defer store.mutex.Unlock()

	if store.users[user.Username] != nil {
		return ErrUserAlreadyExists
	}

	store.users[user.Username] = user.Clone()
	return nil
}

func (store *MapStore) Get(username string) (*model.User, error) {
	store.mutex.RLock()
	defer store.mutex.RUnlock()

	user := store.users[username]
	if user == nil {
		return nil, nil
	}

	return user.Clone(), nil
}
