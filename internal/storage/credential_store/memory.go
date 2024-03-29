package credential_store

import (
	"sync"

	"GophKeeper/internal/server/model/cred"
	"GophKeeper/pkg/errs"
)

type MemoryStorage struct {
	mutex sync.RWMutex
	creds []cred.CredentialFull
}

func NewMemoryStorage() *MemoryStorage {
	return &MemoryStorage{}
}

func (store *MemoryStorage) Create(data cred.CredentialFull) error {

	store.mutex.Lock()
	defer store.mutex.Unlock()

	_, err := store.Find(data.MetaInfo)
	if err == nil {
		return errs.ErrAlreadyExist
	}

	store.creds = append(store.creds, data)
	return nil
}

func (store *MemoryStorage) Get(in cred.CredentialGet) (cred.CredentialFull, error) {

	store.mutex.RLock()
	defer store.mutex.RUnlock()

	idx, err := store.Find(in.MetaInfo)
	if err != nil {
		return cred.CredentialFull{}, err
	}

	return store.creds[idx], nil
}

func (store *MemoryStorage) Delete(in cred.CredentialGet) error {

	store.mutex.Lock()
	defer store.mutex.Unlock()

	idx, err := store.Find(in.MetaInfo)
	if err != nil {
		return err
	}

	// Удаление из найденного элемента из слайса
	store.creds[idx] = store.creds[len(store.creds)-1]
	store.creds = store.creds[:len(store.creds)-1]

	return nil
}

func (store *MemoryStorage) Change(in cred.CredentialFull) error {

	store.mutex.Lock()
	defer store.mutex.Unlock()

	idx, err := store.Find(in.MetaInfo)
	if err != nil {
		return err
	}

	store.creds[idx].Password = in.Password
	return nil
}

func (store *MemoryStorage) Find(metaInfo string) (int, error) {

	for idx, data := range store.creds {
		if data.MetaInfo == metaInfo {
			return idx, nil
		}
	}

	return -1, errs.ErrNotFound
}
