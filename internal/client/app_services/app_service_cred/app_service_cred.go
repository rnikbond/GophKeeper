package app_service_cred

import (
	"GophKeeper/internal/client/model/cred_model"
	"GophKeeper/pkg/errs"
	"GophKeeper/pkg/secret"
	"bufio"
	"crypto/rsa"
	"errors"
	"fmt"
	"github.com/fatih/color"
	"go.uber.org/zap"
	"os"
	"strings"
)

type Sender interface {
	Create(data cred_model.Credential, token string) error
	Get(meta string, token string) (cred_model.Credential, error)
	Delete(meta string, token string) error
	Change(data cred_model.Credential, token string) error
}

type CredOptions func(c *CredService)

type CredService struct {
	Sender

	publicKey  *rsa.PublicKey
	privateKey *rsa.PrivateKey
	logger     *zap.Logger

	token string
}

// NewService - Создание экземпляра сервиса для текстовых данных.
func NewService(s Sender, opts ...CredOptions) *CredService {
	serv := &CredService{
		logger: zap.L(),
		Sender: s,
	}

	for _, opt := range opts {
		opt(serv)
	}

	return serv
}

func WithPublicKey(key *rsa.PublicKey) CredOptions {
	return func(serv *CredService) {
		serv.publicKey = key
	}
}

func WithPrivateKey(key *rsa.PrivateKey) CredOptions {
	return func(serv *CredService) {
		serv.privateKey = key
	}
}

func (serv CredService) ShowMenu() {
	stdin := bufio.NewReader(os.Stdin)

	for {

		fmt.Println("---------------")
		color.Blue(fmt.Sprintf("\tСервис: %s\n", serv.Name()))
		fmt.Println("[0] <- Меню сервисов")
		fmt.Println("[1] Создать")
		fmt.Println("[2] Найти")
		fmt.Println("[3] Удалить")
		fmt.Println("[4] Изменить")
		fmt.Println("---------------")
		fmt.Print("-> ")

		var choice int

		_, err := fmt.Fscan(os.Stdin, &choice)
		stdin.ReadString('\n')
		if err != nil {
			continue
		}

		switch choice {
		case 0:
			return

		case 1:
			serv.Create()

		case 2:
			serv.Get()

		case 3:
			serv.Delete()

		case 4:
			serv.Change()
		}
	}
}

func (serv CredService) Create() {
	data := cred_model.Credential{}

	data.MetaInfo = serv.getInput("Метаинформация: ")
	data.Login = serv.getInputEncode("Логин: ")
	data.Password = serv.getInputEncode("Пароль: ")

	if len(data.MetaInfo) == 0 {
		color.Red("Метаинформация не может быть пустой")
		return
	}

	if len(data.Login) == 0 {
		color.Red("Логин не может быть пустым")
		return
	}

	if len(data.Password) == 0 {
		color.Red("Пароль не может быть пустым")
		return
	}

	err := serv.Sender.Create(data, serv.token)
	if ok := serv.parseError(err); ok {
		color.Green("Данные созданы")
	}
}

func (serv CredService) Get() {
	meta := serv.getInput("Метаинформация: ")

	if len(meta) == 0 {
		color.Red("Метаинформация не может быть пустой")
		return
	}

	data, err := serv.Sender.Get(meta, serv.token)
	if ok := serv.parseError(err); !ok {
		return
	}

	loginDec, errDec := secret.Decrypt(serv.privateKey, data.Login)
	if errDec != nil {
		serv.logger.Error("failed decrypt data", zap.Error(errDec))
		color.Red("Упс... Что-то пошло не так")
		return
	}

	passwordDec, errDec := secret.Decrypt(serv.privateKey, data.Password)
	if errDec != nil {
		serv.logger.Error("failed decrypt data", zap.Error(errDec))
		color.Red("Упс... Что-то пошло не так")
		return
	}

	color.Cyan("Логин : %s", string(loginDec))
	color.Cyan("Пароль: %s", string(passwordDec))
}

func (serv CredService) Delete() {
	meta := serv.getInput("Метаинформация: ")

	if len(meta) == 0 {
		color.Red("Метаинформация не может быть пустой")
		return
	}

	err := serv.Sender.Delete(meta, serv.token)
	if ok := serv.parseError(err); ok {
		color.Green("Данные успешно удалены")
	}
}

func (serv CredService) Change() {
	data := cred_model.Credential{}

	data.MetaInfo = serv.getInput("Метаинформация: ")
	data.Login = serv.getInputEncode("Логин: ")
	data.Password = serv.getInputEncode("Пароль: ")

	if len(data.MetaInfo) == 0 {
		color.Red("Метаинформация не может быть пустой")
		return
	}

	if len(data.Login) == 0 {
		color.Red("Логин не может быть пустым")
		return
	}

	if len(data.Password) == 0 {
		color.Red("Пароль не может быть пустым")
		return
	}

	err := serv.Sender.Change(data, serv.token)
	if ok := serv.parseError(err); ok {
		color.Green("Данные успешно изменены")
	}
}

func (serv CredService) parseError(err error) bool {
	if err == nil {
		return true
	}

	color.New(color.FgRed).Print("\tОшибка: ")

	switch {

	case errors.Is(err, errs.ErrAlreadyExist):
		fmt.Println("Такой метаинформация уже существуют")

	case errors.Is(err, errs.ErrNotFound):
		fmt.Println("Такая метаинформация не найдена")

	case errors.Is(err, errs.ErrLargeData):
		fmt.Println("Размер данных слишком большой")

	default:
		fmt.Println("Внутренняя ошибка сервиса")
		serv.logger.Error("unknown error", zap.Error(err))
	}

	return false
}

func (serv CredService) getInput(title string) string {
	reader := bufio.NewReader(os.Stdin)

	fmt.Print(title)
	data, _ := reader.ReadString('\n')
	data = strings.Replace(data, "\n", "", -1)
	data = strings.Replace(data, "\r", "", -1)

	return data
}

func (serv CredService) getInputEncode(title string) []byte {
	reader := bufio.NewReader(os.Stdin)

	fmt.Print(title)
	data, _ := reader.ReadString('\n')
	data = strings.Replace(data, "\n", "", -1)
	data = strings.Replace(data, "\r", "", -1)

	encodeData, _ := secret.Encrypt(serv.publicKey, []byte(data))
	return encodeData
}

func (serv *CredService) SetToken(token string) {
	serv.token = token
}

func (serv CredService) Name() string {
	return "Логины и пароли"
}
