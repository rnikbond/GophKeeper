//go:generate mockgen -source app_service_credential.go -destination mocks/app_service_credential_mock.go -package app_services
package app_services

import (
	"go.uber.org/zap"

	"GophKeeper/internal/model/cred"
	"GophKeeper/internal/storage/data_store/credential_store"
)

type CredentialApp interface {
	Create(in cred.CredentialFull) error
	Get(in cred.CredentialGet) (cred.CredentialFull, error)
	Delete(in cred.CredentialGet) error
	Change(in cred.CredentialFull) error
}

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
