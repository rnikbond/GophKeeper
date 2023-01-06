package auth_store

import (
	"GophKeeper/internal/model/auth"
	"sync"
)

type MemoryStorage struct {
	mutex sync.RWMutex
	users map[string]auth.Credential
}

func NewMemoryStorage() *MemoryStorage {
	return &MemoryStorage{
		users: make(map[string]auth.Credential),
	}
}

// Create Создание нового пользователя
func (store *MemoryStorage) Create(cred auth.Credential) error {

	store.mutex.Lock()
	defer store.mutex.Unlock()

	if _, ok := store.users[cred.Email]; ok {
		return ErrAlreadyExists
	}

	store.users[cred.Email] = cred
	return nil
}

// Find Поиск пользователя по Email
func (store *MemoryStorage) Find(email string) (auth.Credential, error) {

	store.mutex.RLock()
	defer store.mutex.RUnlock()

	user, ok := store.users[email]
	if !ok {
		return auth.Credential{}, ErrNotFound
	}

	return user, nil
}

// Delete Удаление пользователя
func (store *MemoryStorage) Delete(email string) error {

	store.mutex.Lock()
	defer store.mutex.Unlock()

	if _, ok := store.users[email]; !ok {
		return ErrNotFound
	}

	delete(store.users, email)
	return nil
}

func (store *MemoryStorage) Update(email, password string) error {

	store.mutex.Lock()
	defer store.mutex.Unlock()

	user, ok := store.users[email]
	if !ok {
		return ErrNotFound
	}

	user.Password = password

	store.users[email] = user
	return nil
}
