package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"go.uber.org/zap"

	"GophKeeper/internal/client/clientgrpc"
	"GophKeeper/internal/server"
	"GophKeeper/pkg/logzap"
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

	logzap.ConfigZapLogger()

	logger := zap.L()

	cfg := newConfig()
	cli := newClient(cfg)

	if err := cli.Connect(); err != nil {
		logger.Fatal("connection error", zap.Error(err))
	}

	if err := cli.Login(); err != nil {

		if err = cli.Register(); err != nil {
			logger.Error("register error", zap.Error(err))
		}

		logger.Info("Success Register and Login")
	} else {
		logger.Info("Success Login")
	}

	if err := cli.ChangePassword(); err != nil {
		logger.Error("error change password", zap.Error(err))
	} else {
		logger.Info("success change password")
	}

	if err := cli.CreatePairCred(); err != nil {
		logger.Error("error create credential data", zap.Error(err))
	} else {
		logger.Info("success create credential data")
	}

	if _, err := cli.FindPairCred(); err != nil {
		logger.Error("error find credential data", zap.Error(err))
	} else {
		logger.Info("success find credential data")
	}

	if err := cli.CreateBinary(); err != nil {
		logger.Error("error create binary data", zap.Error(err))
	} else {
		logger.Info("success create binary data")
	}

	if _, err := cli.FindBinary(); err != nil {
		logger.Error("error find binary data", zap.Error(err))
	} else {
		logger.Info("success find binary data")
	}

	if err := cli.CreateText(); err != nil {
		logger.Error("error create text data", zap.Error(err))
	} else {
		logger.Info("success create text data")
	}

	if _, err := cli.FindText(); err != nil {
		logger.Error("error find text data", zap.Error(err))
	} else {
		logger.Info("success find text data")
	}

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
