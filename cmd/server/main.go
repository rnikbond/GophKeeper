package main

import (
	"GophKeeper/internal/storage/auth_store"
	"GophKeeper/internal/storage/data_store/binary_store"
	"GophKeeper/internal/storage/data_store/credential_store"
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

	authStore := auth_store.NewMemoryStorage()
	credStore := credential_store.NewMemoryStorage()
	binStore := binary_store.NewMemoryStorage()

	validate := interceptors.NewValidateInterceptor(cfg.SecretKey)

	authApp := newAuthAppService(authStore, cfg)
	credApp := newCredAppService(credStore)
	binApp := newBinaryAppService(binStore)

	authRPC := newAuthRPCService(authApp)
	credRPC := newCredRPCService(credApp)
	binRPC := newBinaryRPCService(binApp)

	grpcServer, err := servergrpc.NewServer(
		cfg.AddrGRPC,
		validate,
		servergrpc.WithAuthServiceRPC(authRPC),
		servergrpc.WithCredServiceRPC(credRPC),
		servergrpc.WithBinaryServiceRPC(binRPC),
	)

	if err != nil {
		logger := zap.L()
		logger.Error("failed server run", zap.Error(err))
	}

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

func newAuthAppService(store auth_store.AuthStorage, cfg *server.Config) *app_services.AuthAppService {
	return app_services.NewAuthService(store, app_services.WithSecretKey(cfg.SecretKey))
}

func newCredAppService(store credential_store.CredStorage) *app_services.CredentialAppService {
	return app_services.NewCredentialAppService(store)
}

func newBinaryAppService(store binary_store.BinaryStorage) *app_services.BinaryAppService {
	return app_services.NewBinaryAppService(store)
}

func newAuthRPCService(authApp *app_services.AuthAppService) *rpc_services.AuthServiceRPC {
	return rpc_services.NewAuthServiceRPC(authApp)
}

func newCredRPCService(credApp *app_services.CredentialAppService) *rpc_services.CredServiceRPC {
	return rpc_services.NewCredServiceRPC(credApp)
}

func newBinaryRPCService(credApp *app_services.BinaryAppService) *rpc_services.BinaryServiceRPC {
	return rpc_services.NewBinaryServiceRPC(credApp)
}
