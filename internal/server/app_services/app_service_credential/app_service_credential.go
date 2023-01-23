package app_service_credential

import (
	"go.uber.org/zap"

	"GophKeeper/internal/server/model/cred"
	"GophKeeper/internal/storage/credential_store"
)

type CredentialAppService struct {
	store  credential_store.CredStorage
	logger *zap.Logger
}

func NewCredentialAppService(store credential_store.CredStorage) *CredentialAppService {
	return &CredentialAppService{
		store:  store,
		logger: zap.L(),
	}
}

func (serv CredentialAppService) Create(in cred.CredentialFull) error {
	return serv.store.Create(in)
}

func (serv CredentialAppService) Get(in cred.CredentialGet) (cred.CredentialFull, error) {
	return serv.store.Get(in)
}

func (serv CredentialAppService) Delete(in cred.CredentialGet) error {
	return serv.store.Delete(in)
}

func (serv CredentialAppService) Change(in cred.CredentialFull) error {
	return serv.store.Change(in)
}
