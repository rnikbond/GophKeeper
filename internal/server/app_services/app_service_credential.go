//go:generate mockgen -source app_service_credential.go -destination mocks/app_service_credential_mock.go -package app_services
package app_services

type CredentialApp interface {
	Read(email string) (string, error)
	Create(email, password string) error
	Update(email, password string) error
	Delete(email string) error
}

type CredentialAppService struct {
}

func NewCredentialAppService() *CredentialAppService {
	return &CredentialAppService{}
}
