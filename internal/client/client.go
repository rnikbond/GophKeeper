package client

import (
	"GophKeeper/internal/client/app_services/app_service_auth"
	"GophKeeper/pkg/errs"
	"bufio"
	"errors"
	"fmt"
	"github.com/fatih/color"
	"go.uber.org/zap"
	"os"
)

type Options func(c *Client)

type IService interface {
	Name() string
	ShowMenu()
	SetToken(token string)
}

type Client struct {
	logger   *zap.Logger
	auth     *app_service_auth.AuthService
	services []IService
	token    string
}

// NewClient - Создание экземпляра клиента.
func NewClient(auth *app_service_auth.AuthService, opts ...Options) *Client {
	c := &Client{
		auth:   auth,
		logger: zap.L(),
	}

	for _, opt := range opts {
		opt(c)
	}

	return c
}

// WithService - Добавление сервисов клиента.
func WithService(serv IService) Options {
	return func(c *Client) {
		c.services = append(c.services, serv)
	}
}

func (c *Client) Start() {
	token, err := c.auth.Token()

	switch {
	case err == nil:
		break

	case errors.Is(err, errs.ErrCancel):
		color.HiMagenta("Goodbye... :'(")
		return

	default:
		c.logger.Error("failed get token", zap.Error(err))
		color.Red("Увы, но что-то пошло не так...")
		return
	}

	color.Green("Авторизация успешно пройдена")

	for i, _ := range c.services {
		c.services[i].SetToken(token)
	}

	c.showServicesMenu()
	color.HiMagenta("Goodbye... :'(")
}

func (c *Client) showServicesMenu() {

	if len(c.services) == 0 {
		color.Yellow("Досадно, но сервисов нет :(")
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
			c.services[choice-1].ShowMenu()
		}
	}
}
