package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"go.uber.org/zap"

	"GophKeeper/internal/server"
	"GophKeeper/internal/server/app_services"
	"GophKeeper/internal/server/servergrpc"
	"GophKeeper/internal/server/servergrpc/interceptors"
	"GophKeeper/internal/server/servergrpc/rpc_services"
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

	authApp := newAuthAppService(store, cfg)
	authRPC := newAuthRPCService(authApp)
	grpcServer := newServer(cfg, authRPC)

	grpcServer.Start()

	done := make(chan os.Signal)
	signal.Notify(done, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)
	<-done

	grpcServer.Stop()
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

func newAuthAppService(store storage.UserStorage, cfg *server.Config) *app_services.AuthAppService {
	return app_services.NewAuthService(store, app_services.WithSecretKey(cfg.SecretKey))
}

func newAuthRPCService(authApp *app_services.AuthAppService) *rpc_services.AuthServiceRPC {
	return rpc_services.NewAuthServiceRPC(authApp)
}

// newServer Создание объекта сервера
func newServer(cfg *server.Config, auth *rpc_services.AuthServiceRPC) *servergrpc.ServerGRPC {

	validate := interceptors.NewValidateInterceptor(cfg.SecretKey)

	serv, err := servergrpc.NewServer(cfg.AddrGRPC, validate, servergrpc.WithAuthServiceRPC(auth))
	if err != nil {
		logger := zap.L()
		logger.Error("failed server run", zap.Error(err))
	}

	return serv
}
