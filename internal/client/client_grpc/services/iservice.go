package services

type IService interface {
	Name() string
	SetToken(token string)
	ShowMenu() error
}
