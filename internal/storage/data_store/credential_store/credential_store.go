//go:generate mockgen -source credential_store.go -destination mocks/credential_store_mock.go -package credential_store
package credential_store

type CredStorage interface {
	Create() error
	Find() error
	Delete() error
	Update() error
}
