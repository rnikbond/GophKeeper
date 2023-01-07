//go:generate mockgen -source credential_store.go -destination mocks/credential_store_mock.go -package credential_store
package credential_store

import "GophKeeper/internal/model/cred"

type CredStorage interface {
	Create(data cred.CredentialFull) error
	Get(in cred.CredentialGet) (cred.CredentialFull, error)
	Delete(in cred.CredentialGet) error
	Change(in cred.CredentialFull) error
}
