package card_service

import (
	"bufio"
	"context"
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
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"

	"GophKeeper/internal/model/card"
	"GophKeeper/pkg/errs"
	pb "GophKeeper/pkg/proto/card"
	"GophKeeper/pkg/secret"
)

var PeriodLayout = "01.2006"

type CardOptions func(c *CardService)

type CardService struct {
	rpc        pb.CardServiceClient
	logger     *zap.Logger
	publicKey  *rsa.PublicKey
	privateKey *rsa.PrivateKey

	Token string
}

// NewService - Создание экземпляра сервиса для текстовых данных.
func NewService(conn *grpc.ClientConn, opts ...CardOptions) *CardService {

	serv := &CardService{
		rpc:    pb.NewCardServiceClient(conn),
		logger: zap.L(),
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

func (serv CardService) ShowMenu() error {
	if len(serv.Token) == 0 {
		return fmt.Errorf("token is empty")
	}

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
			return nil
		case 1:
			if err := serv.Create(); err != nil {
				if errors.Is(err, errs.ErrAlreadyExist) {
					color.Yellow("Такие данные уже существуют")
				} else if errors.Is(err, ErrInvalidNumber) {
					color.Yellow("Некорректный номер карты")
				} else if errors.Is(err, ErrInvalidPeriod) {
					color.Yellow("Некорректный срок действия карты")
				} else if errors.Is(err, ErrInvalidCVV) {
					color.Yellow("Некорректный CVV")
				} else if errors.Is(err, ErrInvalidFullName) {
					color.Yellow("Некорректные данные держателя карты")
				} else {
					serv.logger.Error("failed create credential data", zap.Error(err))
					color.Red("Внутренняя ошибка при создании данных авторизации")
				}

			} else {
				color.Green("Данные успешно сохранены")
			}

		case 2:

			data, err := serv.Get()
			if err != nil {

				if errors.Is(err, errs.ErrNotFound) {
					color.Red("Такие данные не найдены")
				} else {
					serv.logger.Error("failed delete text data", zap.Error(err))
					color.Red("Внутренняя ошибка при поиске данные")
				}
			} else {
				color.Cyan("Номер карты  : %s", data.Number)
				color.Cyan("Срок действия: %s", data.Period)
				color.Cyan("CVV          : %s", data.CVV)
				color.Cyan("Держатель    : %s", data.FullName)
			}

		case 3:

			if err := serv.Delete(); err != nil {
				if errors.Is(err, errs.ErrNotFound) {
					color.Red("Такие данные не найдены")
				} else {
					serv.logger.Error("failed delete bin data", zap.Error(err))
					color.Red("Внутренняя ошибка при удалении данных")
				}
			} else {
				color.Green("Данные успешно удалены")
			}

		case 4:

			if err := serv.Change(); err != nil {
				if errors.Is(err, errs.ErrNotFound) {
					color.Red("Не найдены данные для изменения")
				} else if errors.Is(err, ErrInvalidNumber) {
					color.Yellow("Некорректный номер карты")
				} else if errors.Is(err, ErrInvalidPeriod) {
					color.Yellow("Некорректный срок действия карты")
				} else if errors.Is(err, ErrInvalidCVV) {
					color.Yellow("Некорректный CVV")
				} else if errors.Is(err, ErrInvalidFullName) {
					color.Yellow("Некорректные данные держателя карты")
				} else {
					serv.logger.Error("failed change bin data", zap.Error(err))
					color.Red("Внутренняя ошибка при изменении данных")
				}
			} else {
				color.Green("Данные успешно изменены")
			}
		}
	}
}

func (serv CardService) Create() error {

	meta := serv.getInput("Метаинформация: ")
	number := serv.getInput("Номер карты: ")
	period := serv.getInput(fmt.Sprintf("Срок действия (формат %s): ", PeriodLayout))
	CVV := serv.getInput("CVV: ")
	fullName := serv.getInput("Держатель: ")

	data := card.DataCardFull{
		Number:   number,
		Period:   period,
		CVV:      CVV,
		FullName: fullName,
	}

	if err := checkCardData(data); err != nil {
		return err
	}

	dataReq := &pb.CreateRequest{
		MetaInfo: meta,
		Number:   serv.encode([]byte(number)),
		Period:   serv.encode([]byte(period)),
		CVV:      serv.encode([]byte(CVV)),
		FullName: serv.encode([]byte(fullName)),
	}

	md := metadata.New(map[string]string{"token": serv.Token})
	ctx := metadata.NewOutgoingContext(context.Background(), md)

	if _, err := serv.rpc.Create(ctx, dataReq); err != nil {
		if e, ok := status.FromError(err); ok {

			switch e.Code() {
			case codes.AlreadyExists:
				return errs.ErrAlreadyExist
			case codes.InvalidArgument:
				return errs.ErrInvalidArgument
			}

			return fmt.Errorf("%d - %s %w", e.Code(), e.String(), errs.ErrInternal)
		}
		return fmt.Errorf("%s %w", err.Error(), errs.ErrInternal)
	}

	return nil
}

func (serv CardService) Get() (card.DataCardFull, error) {

	data := &pb.GetRequest{}
	data.MetaInfo = serv.getInput("Метаинформация: ")

	md := metadata.New(map[string]string{"token": serv.Token})
	ctx := metadata.NewOutgoingContext(context.Background(), md)

	resp, err := serv.rpc.Get(ctx, data)
	if err != nil {
		if e, ok := status.FromError(err); ok {

			switch e.Code() {
			case codes.NotFound:
				return card.DataCardFull{}, errs.ErrNotFound
			}

			return card.DataCardFull{}, fmt.Errorf("%d - %s %w", e.Code(), e.String(), errs.ErrInternal)
		}

		return card.DataCardFull{}, fmt.Errorf("%s %w", err.Error(), errs.ErrInternal)
	}

	numberDecrypt, errDecode := secret.Decrypt(serv.privateKey, resp.Number)
	if errDecode != nil {
		serv.logger.Error("failed decrypt data", zap.Error(errDecode))
		return card.DataCardFull{}, errs.ErrInternal
	}

	periodDecrypt, errDecode := secret.Decrypt(serv.privateKey, resp.Period)
	if errDecode != nil {
		serv.logger.Error("failed decrypt data", zap.Error(errDecode))
		return card.DataCardFull{}, errs.ErrInternal
	}

	cvvDecrypt, errDecode := secret.Decrypt(serv.privateKey, resp.CVV)
	if errDecode != nil {
		serv.logger.Error("failed decrypt data", zap.Error(errDecode))
		return card.DataCardFull{}, errs.ErrInternal
	}

	nameDecrypt, errDecode := secret.Decrypt(serv.privateKey, resp.FullName)
	if errDecode != nil {
		serv.logger.Error("failed decrypt data", zap.Error(errDecode))
		return card.DataCardFull{}, errs.ErrInternal
	}

	return card.DataCardFull{
		Number:   string(numberDecrypt),
		Period:   string(periodDecrypt),
		CVV:      string(cvvDecrypt),
		FullName: string(nameDecrypt),
	}, nil
}

func (serv CardService) Delete() error {

	data := &pb.DeleteRequest{}
	data.MetaInfo = serv.getInput("Метаинформация: ")

	md := metadata.New(map[string]string{"token": serv.Token})
	ctx := metadata.NewOutgoingContext(context.Background(), md)

	_, err := serv.rpc.Delete(ctx, data)
	if err != nil {
		if e, ok := status.FromError(err); ok {

			switch e.Code() {
			case codes.NotFound:
				return errs.ErrNotFound
			}

			return fmt.Errorf("%d - %s %w", e.Code(), e.String(), errs.ErrInternal)
		}

		return fmt.Errorf("%s %w", err.Error(), errs.ErrInternal)
	}

	return nil
}

func (serv CardService) Change() error {

	meta := serv.getInput("Метаинформация: ")
	number := serv.getInput("Номер карты: ")
	period := serv.getInput(fmt.Sprintf("Срок действия (формат %s): ", PeriodLayout))
	CVV := serv.getInput("CVV: ")
	fullName := serv.getInput("Держатель: ")

	data := card.DataCardFull{
		Number:   number,
		Period:   period,
		CVV:      CVV,
		FullName: fullName,
	}

	if err := checkCardData(data); err != nil {
		return err
	}

	dataReq := &pb.ChangeRequest{
		MetaInfo: meta,
		Number:   serv.encode([]byte(number)),
		Period:   serv.encode([]byte(period)),
		CVV:      serv.encode([]byte(CVV)),
		FullName: serv.encode([]byte(fullName)),
	}

	md := metadata.New(map[string]string{"token": serv.Token})
	ctx := metadata.NewOutgoingContext(context.Background(), md)

	_, err := serv.rpc.Change(ctx, dataReq)
	if err != nil {
		if e, ok := status.FromError(err); ok {
			switch e.Code() {
			case codes.NotFound:
				return errs.ErrNotFound
			}

			return fmt.Errorf("%d - %s %w", e.Code(), e.String(), errs.ErrInternal)
		}

		return fmt.Errorf("%s %w", err.Error(), errs.ErrInternal)
	}

	return nil
}

func (serv CardService) getInput(title string) string {

	reader := bufio.NewReader(os.Stdin)

	fmt.Print(title)
	data, _ := reader.ReadString('\n')
	data = strings.Replace(data, "\n", "", -1)
	data = strings.Replace(data, "\r", "", -1)

	return data
}

func (serv CardService) encode(data []byte) []byte {

	encodeData, _ := secret.Encrypt(serv.publicKey, []byte(data))
	return encodeData
}

func (serv CardService) Name() string {
	return "Банковские карты"
}

func (serv *CardService) SetToken(token string) {
	serv.Token = token
}

func checkCardData(in card.DataCardFull) error {

	fmt.Println(in.Period)
	fmt.Println("len", len("in.Period"))
	if _, errTime := time.Parse(PeriodLayout, in.Period); errTime != nil {
		return ErrInvalidPeriod
	}

	if ok, err := luhn.IsValid(in.Number); !ok || err != nil {
		return ErrInvalidNumber
	}

	if len(in.CVV) != 3 {
		return ErrInvalidCVV
	}

	// Используется ParseUint - т.к. не должно быть отрицательного CVV. Например, "-12".
	if _, err := strconv.ParseUint(in.CVV, 10, 32); err != nil {
		return ErrInvalidCVV
	}

	if len(in.FullName) < 4 {
		return ErrInvalidFullName
	}

	return nil
}
