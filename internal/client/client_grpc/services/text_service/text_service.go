package text_service

import (
	"bufio"
	"context"
	"crypto/rsa"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/fatih/color"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"

	"GophKeeper/pkg/errs"
	pb "GophKeeper/pkg/proto/text"
	"GophKeeper/pkg/secret"
)

type TextOptions func(c *TextService)

type TextService struct {
	rpc        pb.TextServiceClient
	logger     *zap.Logger
	publicKey  *rsa.PublicKey
	privateKey *rsa.PrivateKey

	token string
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
	if len(serv.token) == 0 {
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
			err := serv.Create()

			switch {
			case err == nil:
				color.Green("Данные успешно сохранены")

			case errors.Is(err, errs.ErrAlreadyExist):
				color.Yellow("Такие данные уже существуют")

			default:
				serv.logger.Error("failed create text data", zap.Error(err))
				color.Red("Внутренняя ошибка при создании данных")
			}

		case 2:
			text, err := serv.Get()

			switch {
			case err == nil:
				color.Cyan("Данные: %s", text)

			case errors.Is(err, errs.ErrNotFound):
				color.Red("Такие данные не найдены")

			default:
				serv.logger.Error("failed find text data", zap.Error(err))
				color.Red("Внутренняя ошибка при поиске данные")
			}

		case 3:
			err := serv.Delete()

			switch {
			case err == nil:
				color.Green("Данные успешно удалены")

			case errors.Is(err, errs.ErrNotFound):
				color.Red("Такие данные не найдены")

			default:
				serv.logger.Error("failed delete text data", zap.Error(err))
				color.Red("Внутренняя ошибка при удалении данных")
			}

		case 4:
			err := serv.Change()

			switch {
			case err == nil:
				color.Green("Данные успешно изменены")

			case errors.Is(err, errs.ErrNotFound):
				color.Red("Не найдены данные для изменения")

			default:
				serv.logger.Error("failed change text data", zap.Error(err))
				color.Red("Внутренняя ошибка при изменении данных")
			}
		}
	}
}

func (serv TextService) Create() error {

	meta := serv.getInput("Метаинформация: ")
	data := serv.getInputEncode("Текст: ")

	dataReq := &pb.CreateRequest{
		MetaInfo: meta,
		Text:     data,
	}

	md := metadata.New(map[string]string{"token": serv.token})
	ctx := metadata.NewOutgoingContext(context.Background(), md)

	_, err := serv.rpc.Create(ctx, dataReq)
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
	data.MetaInfo = serv.getInput("Метаинформация: ")

	md := metadata.New(map[string]string{"token": serv.token})
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
	data.MetaInfo = serv.getInput("Метаинформация: ")

	md := metadata.New(map[string]string{"token": serv.token})
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

	meta := serv.getInput("Метаинформация: ")
	data := serv.getInputEncode("Текст: ")

	dataReq := &pb.ChangeRequest{
		MetaInfo: meta,
		Text:     data,
	}

	md := metadata.New(map[string]string{"token": serv.token})
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

func (serv TextService) getInput(title string) string {

	reader := bufio.NewReader(os.Stdin)

	fmt.Print(title)
	data, _ := reader.ReadString('\n')
	data = strings.Replace(data, "\n", "", -1)
	data = strings.Replace(data, "\r", "", -1)

	return data
}

func (serv TextService) getInputEncode(title string) []byte {

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

func (serv TextService) Name() string {
	return "Текстовые данные"
}

func (serv *TextService) SetToken(token string) {
	serv.token = token
}
