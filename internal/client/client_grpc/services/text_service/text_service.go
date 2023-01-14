package text_service

import (
	"GophKeeper/pkg/errs"
	pb "GophKeeper/pkg/proto/text"
	"GophKeeper/pkg/secret"
	"bufio"
	"context"
	"crypto/rsa"
	"errors"
	"fmt"
	"github.com/fatih/color"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"os"
)

type TextOptions func(c *TextService)

type TextService struct {
	rpc        pb.TextServiceClient
	logger     *zap.Logger
	publicKey  *rsa.PublicKey
	privateKey *rsa.PrivateKey

	Token string
}

// NewService - Создание экземпляра сервиса для текстовых данных.
func NewService(conn *grpc.ClientConn, opts ...TextOptions) *TextService {

	serv := &TextService{
		rpc:    pb.NewTextServiceClient(conn),
		logger: zap.L(),
	}

	for _, opt := range opts {
		opt(serv)
	}

	return serv
}

func WithPublicKey(key *rsa.PublicKey) TextOptions {
	return func(serv *TextService) {
		serv.publicKey = key
	}
}

func WithPrivateKey(key *rsa.PrivateKey) TextOptions {
	return func(serv *TextService) {
		serv.privateKey = key
	}
}

func (serv TextService) ShowMenu() error {
	if len(serv.Token) == 0 {
		return fmt.Errorf("token is empty")
	}

	stdin := bufio.NewReader(os.Stdin)

	for {

		fmt.Println("---------------")
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
					serv.logger.Error("failed create text data", zap.Error(err))
					color.Red("Внутренняя ошибка при создании текстовых данных")
				}

			} else {
				color.Green("Данные успешно сохранены")
			}

		case 2:

			text, err := serv.Get()
			if err != nil {

				if errors.Is(err, errs.ErrNotFound) {
					color.Red("Такие данные не найдены")
				} else {
					serv.logger.Error("failed delete text data", zap.Error(err))
					color.Red("Внутренняя ошибка при поиске данные")
				}
			} else {
				color.Cyan("Данные: %s", text)
			}

		case 3:

			if err := serv.Delete(); err != nil {
				if errors.Is(err, errs.ErrNotFound) {
					color.Red("Такие данные не найдены")
				} else {
					serv.logger.Error("failed delete text data", zap.Error(err))
					color.Red("Внутренняя ошибка при удалиении данных")
				}
			} else {
				color.Green("Данные успешно удалены")
			}

		case 4:

			if err := serv.Change(); err != nil {
				if errors.Is(err, errs.ErrNotFound) {
					color.Red("Не найдены данные для изменения")
				} else {
					serv.logger.Error("failed change text data", zap.Error(err))
					color.Red("Внутренняя ошибка при изменении данных")
				}
			} else {
				color.Green("Данные успешно изменены")
			}
		}
	}
}

func (serv TextService) Create() error {

	reader := bufio.NewReader(os.Stdin)

	fmt.Print("Метаинформация: ")
	meta, _ := reader.ReadString('\n')

	fmt.Print("Текст: ")
	text, _ := reader.ReadString('\n')

	encodeData, errEncode := secret.Encrypt(serv.publicKey, []byte(text))
	if errEncode != nil {
		serv.logger.Error("failed crypt data", zap.Error(errEncode))
		return errs.ErrInternal
	}

	data := &pb.CreateRequest{
		MetaInfo: meta,
		Text:     encodeData,
	}

	md := metadata.New(map[string]string{"token": serv.Token})
	ctx := metadata.NewOutgoingContext(context.Background(), md)

	_, err := serv.rpc.Create(ctx, data)
	if err != nil {
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

func (serv TextService) Get() (string, error) {

	data := &pb.GetRequest{}
	reader := bufio.NewReader(os.Stdin)

	fmt.Print("Метаинформация: ")
	data.MetaInfo, _ = reader.ReadString('\n')

	md := metadata.New(map[string]string{"token": serv.Token})
	ctx := metadata.NewOutgoingContext(context.Background(), md)

	resp, err := serv.rpc.Get(ctx, data)
	if err != nil {
		if e, ok := status.FromError(err); ok {

			switch e.Code() {
			case codes.NotFound:
				return ``, errs.ErrNotFound
			}

			return ``, fmt.Errorf("%d - %s %w", e.Code(), e.String(), errs.ErrInternal)
		}

		return ``, fmt.Errorf("%s %w", err.Error(), errs.ErrInternal)
	}

	dataDecrypt, errDecode := secret.Decrypt(serv.privateKey, resp.Text)
	if errDecode != nil {
		serv.logger.Error("failed decrypt data", zap.Error(errDecode))
		return ``, errs.ErrInternal
	}

	return string(dataDecrypt), nil
}

func (serv TextService) Delete() error {

	data := &pb.DeleteRequest{}
	reader := bufio.NewReader(os.Stdin)

	fmt.Print("Метаинформация: ")
	data.MetaInfo, _ = reader.ReadString('\n')

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

func (serv TextService) Change() error {

	reader := bufio.NewReader(os.Stdin)

	fmt.Print("Метаинформация: ")
	meta, _ := reader.ReadString('\n')

	fmt.Print("Текст: ")
	text, _ := reader.ReadString('\n')

	encodeData, errEncode := secret.Encrypt(serv.publicKey, []byte(text))
	if errEncode != nil {
		serv.logger.Error("failed crypt data", zap.Error(errEncode))
		return errs.ErrInternal
	}

	data := &pb.ChangeRequest{
		MetaInfo: meta,
		Text:     encodeData,
	}

	md := metadata.New(map[string]string{"token": serv.Token})
	ctx := metadata.NewOutgoingContext(context.Background(), md)

	_, err := serv.rpc.Change(ctx, data)
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

func (serv TextService) Name() string {
	return "Текстовые данные"
}

func (serv *TextService) SetToken(token string) {
	serv.Token = token
}
