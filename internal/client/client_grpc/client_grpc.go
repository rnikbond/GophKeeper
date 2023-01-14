package client_grpc

import (
	"GophKeeper/internal/client/client_grpc/services"
	"GophKeeper/internal/client/client_grpc/services/auth_service"
	"GophKeeper/pkg/errs"
	binary2 "GophKeeper/pkg/proto/binary"
	cardStore "GophKeeper/pkg/proto/card"
	"GophKeeper/pkg/proto/credential"
	textStore "GophKeeper/pkg/proto/text"
	"bufio"
	"errors"
	"fmt"
	"github.com/fatih/color"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"os"
)

type ClientOptions func(c *ClientGRPC)

type ClientGRPC struct {
	conn   *grpc.ClientConn
	logger *zap.Logger

	auth     *auth_service.AuthService
	services []services.IService

	rpcCredClient   credential.CredentialServiceClient
	rpcBinaryClient binary2.BinaryServiceClient
	rpcTextClient   textStore.TextServiceClient
	rpcCardClient   cardStore.CardServiceClient
}

func NewClient(auth *auth_service.AuthService, opts ...ClientOptions) *ClientGRPC {

	c := &ClientGRPC{
		auth:   auth,
		logger: zap.L(),
	}

	for _, opt := range opts {
		opt(c)
	}

	return c
}

func WithService(serv services.IService) ClientOptions {
	return func(c *ClientGRPC) {
		c.services = append(c.services, serv)
	}
}

func (c *ClientGRPC) Start() {

	if !c.showStartMenu() {
		color.HiMagenta(":'(...")
		return
	}

	color.Green("Авторизация успешно пройдена")

	for i, _ := range c.services {
		c.services[i].SetToken(c.auth.Token)
	}

	c.showServicesMenu()
}

func (c *ClientGRPC) showStartMenu() bool {

	stdin := bufio.NewReader(os.Stdin)
	redColor := color.New(color.FgRed)

	for {

		if len(c.auth.Token) != 0 {
			break
		}

		fmt.Println("---------------")
		fmt.Println("[0] Завершить")
		fmt.Println("[1] Авторизация")
		fmt.Println("[2] Регистрация")
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
			return false

		case 1:
			if err := c.auth.SignIn(); err != nil {

				redColor.Print("\tОшибка авторизации: ")

				if errors.Is(err, errs.ErrNotFound) {
					fmt.Println("Пользователь не найден")
				} else if errors.Is(err, errs.ErrInvalidArgument) {
					fmt.Println("Неверный логин или пароль")
				} else {
					fmt.Println("Внутренняя ошибка сервиса")
					c.logger.Fatal("failed sing in", zap.Error(err))
				}
			}

		case 2:
			if err := c.auth.SignUp(); err != nil {

				redColor.Print("\tОшибка регистрации: ")

				if errors.Is(err, errs.ErrAlreadyExist) {
					fmt.Println("Такой Email же зарегистрирован")
				} else if errors.Is(err, errs.ErrInvalidArgument) {
					fmt.Println("Некорректный Email или пароль слишком короткий")
				} else {
					fmt.Println("Внутренняя ошибка сервиса")
					//c.logger.Fatal("failed sing up", zap.Error(err))
				}
			}
		}
	}

	return true
}

func (c *ClientGRPC) showServicesMenu() {

	if len(c.services) == 0 {
		color.Yellow("Извините, сервисов нет :(")
		return
	}

	stdin := bufio.NewReader(os.Stdin)

	for {

		fmt.Println("---------------")
		fmt.Println("[0] Завершить")
		for i, serv := range c.services {
			fmt.Printf("[%d] %s\n", i+1, serv.Name())
		}
		fmt.Println("---------------")
		fmt.Print("-> ")

		var choice int

		_, err := fmt.Fscan(os.Stdin, &choice)
		stdin.ReadString('\n')

		if err != nil {
			continue
		}

		if choice == 0 {
			return
		}

		if choice >= 1 && choice <= len(c.services) {
			if err := c.services[choice-1].ShowMenu(); err != nil {
				c.logger.Error(fmt.Sprintf("failed run menu service %s", c.services[choice-1].Name()), zap.Error(err))
			}
		}
	}
}
