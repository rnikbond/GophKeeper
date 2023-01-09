package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"

	"go.uber.org/zap"

	"GophKeeper/internal/server"
	"GophKeeper/internal/server/app_services"
	"GophKeeper/internal/server/servergrpc"
	"GophKeeper/internal/server/servergrpc/interceptors"
	"GophKeeper/internal/server/servergrpc/rpc_services"
	"GophKeeper/internal/storage/auth_store"
	"GophKeeper/internal/storage/data_store/binary_store"
	"GophKeeper/internal/storage/data_store/card_store"
	"GophKeeper/internal/storage/data_store/credential_store"
	"GophKeeper/internal/storage/data_store/text_store"
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

	logger := zap.L()
	cfg := newConfig()

	var authStore auth_store.AuthStorage
	var credStore credential_store.CredStorage
	var binStore binary_store.BinaryStorage
	var textStore text_store.TextStorage
	var cardStore card_store.CardStorage

	// Создание хранилищ
	if len(cfg.DatabaseURI) != 0 {

		db, errDB := initDB(cfg.DatabaseURI)
		if errDB != nil {
			logger.Fatal("failed database connect", zap.Error(errDB))
		}

		authStore = auth_store.NewPostgresStorage(db)

		credStore = credential_store.NewMemoryStorage()
		binStore = binary_store.NewMemoryStorage()
		textStore = text_store.NewMemoryStorage()
		cardStore = card_store.NewMemoryStorage()

	} else {
		authStore = auth_store.NewMemoryStorage()
		credStore = credential_store.NewMemoryStorage()
		binStore = binary_store.NewMemoryStorage()
		textStore = text_store.NewMemoryStorage()
		cardStore = card_store.NewMemoryStorage()
	}

	// Создание сервисов приложения
	authApp := app_services.NewAuthService(authStore, app_services.WithSecretKey(cfg.SecretKey))
	credApp := app_services.NewCredentialAppService(credStore)
	binApp := app_services.NewBinaryAppService(binStore)
	textApp := app_services.NewTextAppService(textStore)
	cardApp := app_services.NewCardAppService(cardStore)

	// Создание gRPC сервисов
	authRPC := rpc_services.NewAuthServiceRPC(authApp)
	credRPC := rpc_services.NewCredServiceRPC(credApp)
	binRPC := rpc_services.NewBinaryServiceRPC(binApp)
	textRPC := rpc_services.NewTextServiceRPC(textApp)
	cardRPC := rpc_services.NewCardServiceRPC(cardApp)

	validate := interceptors.NewValidateInterceptor(cfg.SecretKey)

	// Создание сервера
	grpcServer, err := servergrpc.NewServer(
		cfg.AddrGRPC,
		validate,
		servergrpc.WithAuthServiceRPC(authRPC),
		servergrpc.WithCredServiceRPC(credRPC),
		servergrpc.WithBinaryServiceRPC(binRPC),
		servergrpc.WithTextServiceRPC(textRPC),
		servergrpc.WithCardServiceRPC(cardRPC),
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
