package main

import (
	"GophKeeper/internal/client/clientgrpc"
	"GophKeeper/internal/server"
	"GophKeeper/pkg/logzap"
	"fmt"
	"go.uber.org/zap"
	"os"
	"os/signal"
	"syscall"
)

var (
	_ = (*clientgrpc.ClientGRPC)(nil)
)

var (
	buildVersion = "N/A"
	buildDate    = "N/A"
	buildCommit  = "N/A"
)

func main() {

	logzap.ConfigZapLogger("client_errs.log")

	logger := zap.L()

	cfg := newConfig()
	cli := newClient(cfg)

	if err := cli.Connect(); err != nil {
		logger.Fatal(" connection error", zap.Error(err))
	}

	if err := cli.Login(); err != nil {

		if err = cli.Register(); err != nil {
			logger.Error("register error", zap.Error(err))
		}

		logger.Info("Success Register")
	} else {
		logger.Error("login error", zap.Error(err))
	}

	logger.Info("Success Login")

	done := make(chan os.Signal)
	signal.Notify(done, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)
	<-done

	if err := cli.Disconnect(); err != nil {
		logger.Error("could not disconnect", zap.Error(err))
	}

	logger.Info("Goodbye...")
}

func init() {

	fmt.Printf("Build version: %s\n", buildVersion)
	fmt.Printf("Build date: %s\n", buildDate)
	fmt.Printf("Build commit: %s\n", buildCommit)
}

// newServer Создание объекта конфигурации
func newConfig() *server.Config {
	cfg := server.NewConfig()
	if err := cfg.ParseArgs(); err != nil {
		logger := zap.L()
		logger.Fatal("failed run gRPC server: %v\n", zap.Error(err))
	}

	return cfg
}

// newClient Создание объекта клиента
func newClient(cfg *server.Config) *clientgrpc.ClientGRPC {

	return clientgrpc.NewClient(cfg.AddrGRPC)
}
