package binary_store

import (
	"sync"

	"GophKeeper/internal/server/model/binary"
	"GophKeeper/pkg/errs"
)

type MemoryStorage struct {
	mutex sync.RWMutex
	creds []binary.DataFull
}

func NewMemoryStorage() *MemoryStorage {
	return &MemoryStorage{}
}

func (store *MemoryStorage) Create(in binary.DataFull) error {

	store.mutex.Lock()
	defer store.mutex.Unlock()

	_, err := store.Find(in.MetaInfo)
	if err == nil {
		return errs.ErrAlreadyExist
	}

	store.creds = append(store.creds, in)
	return nil
}

func (store *MemoryStorage) Get(in binary.DataGet) (binary.DataFull, error) {

	store.mutex.RLock()
	defer store.mutex.RUnlock()

	idx, err := store.Find(in.MetaInfo)
	if err != nil {
		return binary.DataFull{}, err
	}

	return store.creds[idx], nil
}

func (store *MemoryStorage) Delete(in binary.DataGet) error {

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

func (store *MemoryStorage) Change(in binary.DataFull) error {

	store.mutex.Lock()
	defer store.mutex.Unlock()

	idx, err := store.Find(in.MetaInfo)
	if err != nil {
		return err
	}

	store.creds[idx].Bytes = in.Bytes
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
