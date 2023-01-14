package binary_service

import (
	"GophKeeper/pkg/errs"
	pb "GophKeeper/pkg/proto/binary"
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
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

type BinaryOptions func(c *BinaryService)

type BinaryService struct {
	rpc        pb.BinaryServiceClient
	logger     *zap.Logger
	publicKey  *rsa.PublicKey
	privateKey *rsa.PrivateKey

	Token string
}

// NewService - Создание экземпляра сервиса для текстовых данных.
func NewService(conn *grpc.ClientConn, opts ...BinaryOptions) *BinaryService {

	serv := &BinaryService{
		rpc:    pb.NewBinaryServiceClient(conn),
		logger: zap.L(),
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

func (serv BinaryService) ShowMenu() error {
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
				} else if errors.Is(err, errs.ErrNotFound) {
					color.Yellow("Файл с бинарными данными не найден")
				} else {
					serv.logger.Error("failed create text data", zap.Error(err))
					color.Red("Внутренняя ошибка при создании бинарных данных")
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
					serv.logger.Error("failed delete bin data", zap.Error(err))
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
					serv.logger.Error("failed change bin data", zap.Error(err))
					color.Red("Внутренняя ошибка при изменении данных")
				}
			} else {
				color.Green("Данные успешно изменены")
			}
		}
	}
}

func (serv BinaryService) Create() error {

	ex, err := os.Executable()
	if err != nil {
		panic(err)
	}
	exPath := filepath.Dir(ex)
	fmt.Println(exPath)

	reader := bufio.NewReader(os.Stdin)

	fmt.Print("Метаинформация: ")
	meta, _ := reader.ReadString('\n')
	meta = strings.Replace(meta, "\n", "", -1)
	meta = strings.Replace(meta, "\r", "", -1)

	fmt.Print("Данные: ")
	data, _ := reader.ReadString('\n')
	data = strings.Replace(data, "\n", "", -1)
	data = strings.Replace(data, "\r", "", -1)

	// Выбрали путь к файлу
	if _, err := os.Stat(data); err == nil {
		fileData, errRead := ioutil.ReadFile(data)
		if errRead != nil {
			return errs.ErrNotFound
		}

		color.Cyan("Выбран файл")
		data = string(fileData)
	} else {
		color.Cyan("Вы ввели данные вручную")
	}

	encodeData, errEncode := secret.Encrypt(serv.publicKey, []byte(data))
	if errEncode != nil {
		serv.logger.Error("failed crypt data", zap.Error(errEncode))
		return errs.ErrInternal
	}

	binData := &pb.CreateRequest{
		MetaInfo: meta,
		Data:     encodeData,
	}

	md := metadata.New(map[string]string{"token": serv.Token})
	ctx := metadata.NewOutgoingContext(context.Background(), md)

	_, err = serv.rpc.Create(ctx, binData)
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

func (serv BinaryService) Get() (string, error) {

	data := &pb.GetRequest{}
	reader := bufio.NewReader(os.Stdin)

	fmt.Print("Метаинформация: ")
	data.MetaInfo, _ = reader.ReadString('\n')
	data.MetaInfo = strings.Replace(data.MetaInfo, "\n", "", -1)
	data.MetaInfo = strings.Replace(data.MetaInfo, "\r", "", -1)

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

	dataDecrypt, errDecode := secret.Decrypt(serv.privateKey, resp.Data)
	if errDecode != nil {
		serv.logger.Error("failed decrypt data", zap.Error(errDecode))
		return ``, errs.ErrInternal
	}

	return string(dataDecrypt), nil
}

func (serv BinaryService) Delete() error {

	data := &pb.DeleteRequest{}
	reader := bufio.NewReader(os.Stdin)

	fmt.Print("Метаинформация: ")
	data.MetaInfo, _ = reader.ReadString('\n')
	data.MetaInfo = strings.Replace(data.MetaInfo, "\n", "", -1)
	data.MetaInfo = strings.Replace(data.MetaInfo, "\r", "", -1)

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

func (serv BinaryService) Change() error {

	reader := bufio.NewReader(os.Stdin)

	fmt.Print("Метаинформация: ")
	meta, _ := reader.ReadString('\n')
	meta = strings.Replace(meta, "\n", "", -1)
	meta = strings.Replace(meta, "\r", "", -1)

	fmt.Print("Текст: ")
	data, _ := reader.ReadString('\n')
	data = strings.Replace(data, "\n", "", -1)
	data = strings.Replace(data, "\r", "", -1)

	// Выбрали путь к файлу
	if _, err := os.Stat(data); err == nil {
		fileData, errRead := ioutil.ReadFile(data)
		if errRead != nil {
			return errs.ErrNotFound
		}

		color.Cyan("Выбран файл")
		data = string(fileData)
	} else {
		color.Cyan("Вы ввели данные вручную")
	}

	encodeData, errEncode := secret.Encrypt(serv.publicKey, []byte(data))
	if errEncode != nil {
		serv.logger.Error("failed crypt data", zap.Error(errEncode))
		return errs.ErrInternal
	}

	binData := &pb.ChangeRequest{
		MetaInfo: meta,
		Data:     encodeData,
	}

	md := metadata.New(map[string]string{"token": serv.Token})
	ctx := metadata.NewOutgoingContext(context.Background(), md)

	_, err := serv.rpc.Change(ctx, binData)
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

func (serv BinaryService) Name() string {
	return "Бинарные данные"
}

func (serv *BinaryService) SetToken(token string) {
	serv.Token = token
}
