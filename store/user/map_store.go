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
		return nil, ErrUserNotFound
	}

	return user.Clone(), nil
}

func (store *MapStore) ExistsByEmail(email string) bool {
	store.mutex.RLock()
	defer store.mutex.RUnlock()

	for _, v := range store.users {
		if v.Email == email {
			return true
		}
	}
	return false
}

func (store *MapStore) ExistsByUsername(username string) bool {
	store.mutex.RLock()
	defer store.mutex.RUnlock()

	_, ok := store.users[username]
	return ok
}
