package app_service_binary

import (
	"bufio"
	"crypto/rsa"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/fatih/color"
	"go.uber.org/zap"

	"GophKeeper/internal/client/model/binary_model"
	"GophKeeper/pkg/errs"
	"GophKeeper/pkg/secret"
)

type Sender interface {
	Create(text binary_model.Binary, token string) error
	Get(meta string, token string) (binary_model.Binary, error)
	Delete(meta string, token string) error
	Change(text binary_model.Binary, token string) error
}

type BinaryOptions func(c *BinaryService)

type BinaryService struct {
	Sender

	publicKey  *rsa.PublicKey
	privateKey *rsa.PrivateKey
	logger     *zap.Logger

	token string
}

// NewService - Создание экземпляра сервиса для текстовых данных.
func NewService(s Sender, opts ...BinaryOptions) *BinaryService {
	serv := &BinaryService{
		logger: zap.L(),
		Sender: s,
	}

	for _, opt := range opts {
		opt(serv)
	}

	return serv
}

func WithPublicKey(key *rsa.PublicKey) BinaryOptions {
	return func(serv *BinaryService) {
		serv.publicKey = key
	}
}

func WithPrivateKey(key *rsa.PrivateKey) BinaryOptions {
	return func(serv *BinaryService) {
		serv.privateKey = key
	}
}

func (serv BinaryService) ShowMenu() {

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

func (serv BinaryService) Create() {

	data := binary_model.Binary{}

	data.MetaInfo = serv.getInput("Метаинформация: ")
	data.Data = serv.getInputEncode("Данные: ")

	if len(data.MetaInfo) == 0 {
		color.Red("Метаинформация не может быть пустой")
		return
	}

	if len(data.Data) == 0 {
		color.Red("Данные не могут быть пустыми")
		return
	}

	err := serv.Sender.Create(data, serv.token)
	if ok := serv.parseError(err); ok {
		color.Green("Данные созданы")
	}
}

func (serv BinaryService) Get() {

	meta := serv.getInput("Метаинформация: ")

	if len(meta) == 0 {
		color.Red("Метаинформация не может быть пустой")
		return
	}

	text, err := serv.Sender.Get(meta, serv.token)
	if ok := serv.parseError(err); !ok {
		return
	}

	dataDec, errDec := secret.Decrypt(serv.privateKey, text.Data)
	if errDec != nil {
		serv.logger.Error("failed decrypt data", zap.Error(errDec))
		color.Red("Упс... Что-то пошло не так")
		return
	}

	color.Cyan("Данные: %s", string(dataDec))
}

func (serv BinaryService) Delete() {

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

func (serv BinaryService) Change() {

	data := binary_model.Binary{}

	data.MetaInfo = serv.getInput("Метаинформация: ")
	data.Data = serv.getInputEncode("Текст: ")

	if len(data.MetaInfo) == 0 {
		color.Red("Метаинформация не может быть пустой")
		return
	}

	if len(data.Data) == 0 {
		color.Red("Данные не могут быть пустыми")
		return
	}

	err := serv.Sender.Change(data, serv.token)
	if ok := serv.parseError(err); ok {
		color.Green("Данные успешно изменены")
	}
}

func (serv BinaryService) parseError(err error) bool {

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

func (serv BinaryService) getInput(title string) string {

	reader := bufio.NewReader(os.Stdin)

	fmt.Print(title)
	data, _ := reader.ReadString('\n')
	data = strings.Replace(data, "\n", "", -1)
	data = strings.Replace(data, "\r", "", -1)

	return data
}

func (serv BinaryService) getInputEncode(title string) []byte {

	reader := bufio.NewReader(os.Stdin)

	fmt.Print(title)
	data, _ := reader.ReadString('\n')
	data = strings.Replace(data, "\n", "", -1)
	data = strings.Replace(data, "\r", "", -1)

	// Выбрали путь к файлу
	if _, err := os.Stat(data); err == nil {
		fileData, errRead := ioutil.ReadFile(data)
		if errRead != nil {
			return nil
		}

		color.Cyan("Выбран файл")
		data = string(fileData)
	} else {
		color.Cyan("Вы ввели данные вручную")
	}

	encodeData, _ := secret.Encrypt(serv.publicKey, []byte(data))
	return encodeData
}

func (serv *BinaryService) SetToken(token string) {
	serv.token = token
}

func (serv BinaryService) Name() string {
	return "Бинарные данные"
}
