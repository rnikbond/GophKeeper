package cred_service

import (
	"bufio"
	"context"
	"crypto/rsa"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/fatih/color"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"

	"GophKeeper/pkg/errs"
	pb "GophKeeper/pkg/proto/credential"
	"GophKeeper/pkg/secret"
)

type CredOptions func(c *CredService)

type CredService struct {
	rpc        pb.CredentialServiceClient
	logger     *zap.Logger
	publicKey  *rsa.PublicKey
	privateKey *rsa.PrivateKey

	Token string
}

// NewService - Создание экземпляра сервиса для текстовых данных.
func NewService(conn *grpc.ClientConn, opts ...CredOptions) *CredService {

	serv := &CredService{
		rpc:    pb.NewCredentialServiceClient(conn),
		logger: zap.L(),
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

func (serv CredService) ShowMenu() error {
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
				} else {
					serv.logger.Error("failed create credential data", zap.Error(err))
					color.Red("Внутренняя ошибка при создании данных авторизации")
				}

			} else {
				color.Green("Данные успешно сохранены")
			}

		case 2:

			login, pwd, err := serv.Get()
			if err != nil {

				if errors.Is(err, errs.ErrNotFound) {
					color.Red("Такие данные не найдены")
				} else {
					serv.logger.Error("failed delete text data", zap.Error(err))
					color.Red("Внутренняя ошибка при поиске данные")
				}
			} else {
				color.Cyan("Логин : %s", login)
				color.Cyan("Пароль: %s", pwd)
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

func (serv CredService) Create() error {

	meta := serv.getInput("Метаинформация: ")
	login := serv.getInputEncode("Логин: ")
	password := serv.getInputEncode("Пароль: ")

	dataReq := &pb.CreateRequest{
		MetaInfo: meta,
		Email:    login,
		Password: password,
	}

	md := metadata.New(map[string]string{"token": serv.Token})
	ctx := metadata.NewOutgoingContext(context.Background(), md)

	if _, err := serv.rpc.Create(ctx, dataReq); err != nil {
		if e, ok := status.FromError(err); ok {

			switch e.Code() {
			case codes.AlreadyExists:
				return errs.ErrAlreadyExist
			}

			return fmt.Errorf("%d - %s %w", e.Code(), e.String(), errs.ErrInternal)
		}
		return fmt.Errorf("%s %w", err.Error(), errs.ErrInternal)
	}

	return nil
}

func (serv CredService) Get() (string, string, error) {

	data := &pb.GetRequest{}
	data.MetaInfo = serv.getInput("Метаинформация: ")

	md := metadata.New(map[string]string{"token": serv.Token})
	ctx := metadata.NewOutgoingContext(context.Background(), md)

	resp, err := serv.rpc.Get(ctx, data)
	if err != nil {
		if e, ok := status.FromError(err); ok {

			switch e.Code() {
			case codes.NotFound:
				return ``, ``, errs.ErrNotFound
			}

			return ``, ``, fmt.Errorf("%d - %s %w", e.Code(), e.String(), errs.ErrInternal)
		}

		return ``, ``, fmt.Errorf("%s %w", err.Error(), errs.ErrInternal)
	}

	loginDecrypt, errDecode := secret.Decrypt(serv.privateKey, resp.Email)
	if errDecode != nil {
		serv.logger.Error("failed decrypt data", zap.Error(errDecode))
		return ``, ``, errs.ErrInternal
	}

	passwordDecrypt, errDecode := secret.Decrypt(serv.privateKey, resp.Password)
	if errDecode != nil {
		serv.logger.Error("failed decrypt data", zap.Error(errDecode))
		return ``, ``, errs.ErrInternal
	}

	return string(loginDecrypt), string(passwordDecrypt), nil
}

func (serv CredService) Delete() error {

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

func (serv CredService) Change() error {

	meta := serv.getInput("Метаинформация: ")
	login := serv.getInputEncode("Логин: ")
	password := serv.getInputEncode("Пароль: ")

	dataReq := &pb.ChangeRequest{
		MetaInfo: meta,
		Email:    login,
		Password: password,
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

func (serv CredService) Name() string {
	return "Логины и пароли"
}

func (serv *CredService) SetToken(token string) {
	serv.Token = token
}
