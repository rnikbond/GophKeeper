package text_store

import (
	"sync"

	"GophKeeper/internal/model/text"
	"GophKeeper/pkg/errs"
)

type MemoryStorage struct {
	mutex sync.RWMutex
	data  []text.DataTextFull
}

func NewMemoryStorage() *MemoryStorage {
	return &MemoryStorage{}
}

func (store *MemoryStorage) Create(data text.DataTextFull) error {

	store.mutex.Lock()
	defer store.mutex.Unlock()

	_, err := store.Find(data.MetaInfo)
	if err == nil {
		return errs.ErrAlreadyExist
	}

	store.data = append(store.data, data)
	return nil
}

func (store *MemoryStorage) Get(in text.DataTextGet) (text.DataTextFull, error) {

	store.mutex.RLock()
	defer store.mutex.RUnlock()

	idx, err := store.Find(in.MetaInfo)
	if err != nil {
		return text.DataTextFull{}, err
	}

	return store.data[idx], nil
}

func (store *MemoryStorage) Delete(in text.DataTextGet) error {

	store.mutex.Lock()
	defer store.mutex.Unlock()

	idx, err := store.Find(in.MetaInfo)
	if err != nil {
		return err
	}

	// Удаление из найденного элемента из слайса
	store.data[idx] = store.data[len(store.data)-1]
	store.data = store.data[:len(store.data)-1]

	return nil
}

func (store *MemoryStorage) Change(in text.DataTextFull) error {

	store.mutex.Lock()
	defer store.mutex.Unlock()

	idx, err := store.Find(in.MetaInfo)
	if err != nil {
		return err
	}

	store.data[idx].Text = in.Text
	return nil
}

func (store *MemoryStorage) Find(metaInfo string) (int, error) {

	for idx, data := range store.data {
		if data.MetaInfo == metaInfo {
			return idx, nil
		}
	}

	return -1, errs.ErrNotFound
}
