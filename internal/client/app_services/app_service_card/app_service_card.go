package app_service_card

import (
	"bufio"
	"crypto/rsa"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/EClaesson/go-luhn"
	"github.com/fatih/color"
	"go.uber.org/zap"

	"GophKeeper/internal/client/model/card_model"
	"GophKeeper/pkg/errs"
	"GophKeeper/pkg/secret"
)

var PeriodLayout = "01.2006"

type Sender interface {
	Create(data card_model.Card, token string) error
	Get(meta string, token string) (card_model.Card, error)
	Delete(meta string, token string) error
	Change(data card_model.Card, token string) error
}

type CardOptions func(c *CardService)

type CardService struct {
	Sender

	publicKey  *rsa.PublicKey
	privateKey *rsa.PrivateKey
	logger     *zap.Logger

	token string
}

// NewService - Создание экземпляра сервиса для текстовых данных.
func NewService(s Sender, opts ...CardOptions) *CardService {
	serv := &CardService{
		logger: zap.L(),
		Sender: s,
	}

	for _, opt := range opts {
		opt(serv)
	}

	return serv
}

func WithPublicKey(key *rsa.PublicKey) CardOptions {
	return func(serv *CardService) {
		serv.publicKey = key
	}
}

func WithPrivateKey(key *rsa.PrivateKey) CardOptions {
	return func(serv *CardService) {
		serv.privateKey = key
	}
}

func (serv CardService) ShowMenu() {
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

func (serv CardService) Create() {
	data := card_model.Card{}

	data.MetaInfo = serv.getInput("Метаинформация: ")
	number := serv.getInput("Номер: ")
	period := serv.getInput("Период: ")
	CVV := serv.getInput("CVV: ")
	holder := serv.getInput("Держатель: ")

	if len(data.MetaInfo) == 0 {
		color.Red("Метаинформация не может быть пустой")
		return
	}

	if ok := serv.checkCardData(number, period, CVV, holder); !ok {
		return
	}

	data.Number = serv.encode(number)
	data.Period = serv.encode(period)
	data.CVV = serv.encode(CVV)
	data.FullName = serv.encode(holder)

	err := serv.Sender.Create(data, serv.token)
	if ok := serv.parseError(err); ok {
		color.Green("Данные созданы")
	}
}

func (serv CardService) Get() {
	meta := serv.getInput("Метаинформация: ")

	if len(meta) == 0 {
		color.Red("Метаинформация не может быть пустой")
		return
	}

	data, err := serv.Sender.Get(meta, serv.token)
	if ok := serv.parseError(err); !ok {
		return
	}

	number, errDec := secret.Decrypt(serv.privateKey, data.Number)
	if errDec != nil {
		serv.logger.Error("failed decrypt data", zap.Error(errDec))
		color.Red("Упс... Что-то пошло не так")
		return
	}

	period, errDec := secret.Decrypt(serv.privateKey, data.Period)
	if errDec != nil {
		serv.logger.Error("failed decrypt data", zap.Error(errDec))
		color.Red("Упс... Что-то пошло не так")
		return
	}

	CVV, errDec := secret.Decrypt(serv.privateKey, data.CVV)
	if errDec != nil {
		serv.logger.Error("failed decrypt data", zap.Error(errDec))
		color.Red("Упс... Что-то пошло не так")
		return
	}

	holder, errDec := secret.Decrypt(serv.privateKey, data.FullName)
	if errDec != nil {
		serv.logger.Error("failed decrypt data", zap.Error(errDec))
		color.Red("Упс... Что-то пошло не так")
		return
	}

	color.Cyan("Номер    : %s", string(number))
	color.Cyan("Период   : %s", string(period))
	color.Cyan("CVV      : %s", string(CVV))
	color.Cyan("Держатель: %s", string(holder))
}

func (serv CardService) Delete() {
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

func (serv CardService) Change() {
	data := card_model.Card{}

	data.MetaInfo = serv.getInput("Метаинформация: ")
	number := serv.getInput("Номер: ")
	period := serv.getInput("Период: ")
	CVV := serv.getInput("CVV: ")
	holder := serv.getInput("Держатель: ")

	if len(data.MetaInfo) == 0 {
		color.Red("Метаинформация не может быть пустой")
		return
	}

	if ok := serv.checkCardData(number, period, CVV, holder); !ok {
		return
	}

	data.Number = serv.encode(number)
	data.Period = serv.encode(period)
	data.CVV = serv.encode(CVV)
	data.FullName = serv.encode(holder)

	err := serv.Sender.Change(data, serv.token)
	if ok := serv.parseError(err); ok {
		color.Green("Данные успешно изменены")
	}
}

func (serv CardService) parseError(err error) bool {
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

func (serv CardService) getInput(title string) string {
	reader := bufio.NewReader(os.Stdin)

	fmt.Print(title)
	data, _ := reader.ReadString('\n')
	data = strings.Replace(data, "\n", "", -1)
	data = strings.Replace(data, "\r", "", -1)

	return data
}

func (serv CardService) encode(data string) []byte {

	encodeData, _ := secret.Encrypt(serv.publicKey, []byte(data))
	return encodeData
}

func (serv CardService) checkCardData(number, period, cvv, holder string) bool {

	if _, errTime := time.Parse(PeriodLayout, period); errTime != nil {
		color.Red("Некорректный период")
		return false
	}

	if len(number) != 16 {
		color.Red("Некорректный номер карты")
		return false
	}

	if ok, err := luhn.IsValid(number); !ok || err != nil {
		color.Red("Некорректный номер карты")
		return false
	}

	if len(cvv) != 3 {
		color.Red("Некорректный CVV")
		return false
	}

	// Используется ParseUint - т.к. не должно быть отрицательного CVV. Например, "-12".
	if _, err := strconv.ParseUint(cvv, 10, 32); err != nil {
		color.Red("Некорректный CVV")
		return false
	}

	if len(holder) < 4 {
		color.Red("Некорректные данные держателя")
		return false
	}

	return true
}

func (serv *CardService) SetToken(token string) {
	serv.token = token
}

func (serv CardService) Name() string {
	return "Банковские карты"
}
