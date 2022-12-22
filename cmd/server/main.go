package main

import (
	"GophKeeper/internal/server/model"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"go.uber.org/zap"

	"GophKeeper/internal/server"
	"GophKeeper/internal/server/servergrpc"
	"GophKeeper/internal/storage"
	"GophKeeper/pkg/logzap"
)

var (
	_ = (*servergrpc.ServerGRPC)(nil)
)

var (
	buildVersion = "N/A"
	buildDate    = "N/A"
	buildCommit  = "N/A"
)

func main() {

	logzap.ConfigZapLogger()

	cfg := newConfig()
	store := storage.NewMemoryStorage()
	auth := newAuthModel(store)
	serv := newServer(cfg, auth)

	serv.Start()

	done := make(chan os.Signal)
	signal.Notify(done, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)
	<-done

	serv.Stop()
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
		logger.Error("failed parse args", zap.Error(err))
	}

	return cfg
}

func newAuthModel(store storage.UserStorage) *model.AuthModel {
	return model.NewAuthModel(store)
}

// newServer Создание объекта сервера
func newServer(cfg *server.Config, auth *model.AuthModel) *servergrpc.ServerGRPC {

	serv, err := servergrpc.NewServer(cfg.AddrGRPC, auth)
	if err != nil {
		logger := zap.L()
		logger.Error("failed server run", zap.Error(err))
	}

	return serv
}
