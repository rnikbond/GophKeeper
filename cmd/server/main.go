package main

import (
	"GophKeeper/internal/server/app_services/app_service_auth"
	"GophKeeper/internal/server/app_services/app_service_binary"
	"GophKeeper/internal/server/app_services/app_service_card"
	"GophKeeper/internal/server/app_services/app_service_credential"
	"GophKeeper/internal/server/app_services/app_service_text"
	binary_store2 "GophKeeper/internal/storage/binary_store"
	card_store2 "GophKeeper/internal/storage/card_store"
	credential_store2 "GophKeeper/internal/storage/credential_store"
	text_store2 "GophKeeper/internal/storage/text_store"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"

	"go.uber.org/zap"

	"GophKeeper/internal/server"
	"GophKeeper/internal/server/servergrpc"
	"GophKeeper/internal/server/servergrpc/interceptors"
	"GophKeeper/internal/server/servergrpc/rpc_services"
	"GophKeeper/internal/storage/auth_store"
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
	var credStore credential_store2.CredStorage
	var binStore binary_store2.BinaryStorage
	var textStore text_store2.TextStorage
	var cardStore card_store2.CardStorage

	// Создание хранилищ
	if len(cfg.DatabaseURI) != 0 {

		db, errDB := initDB(cfg.DatabaseURI)
		if errDB != nil {
			logger.Fatal("failed database connect", zap.Error(errDB))
		}

		authStore = auth_store.NewPostgresStorage(db)

		credStore = credential_store2.NewMemoryStorage()
		binStore = binary_store2.NewMemoryStorage()
		textStore = text_store2.NewMemoryStorage()
		cardStore = card_store2.NewMemoryStorage()

	} else {
		authStore = auth_store.NewMemoryStorage()
		credStore = credential_store2.NewMemoryStorage()
		binStore = binary_store2.NewMemoryStorage()
		textStore = text_store2.NewMemoryStorage()
		cardStore = card_store2.NewMemoryStorage()
	}

	// Создание сервисов приложения
	authApp := app_service_auth.NewAuthService(authStore, app_service_auth.WithSecretKey(cfg.SecretKey))
	credApp := app_service_credential.NewCredentialAppService(credStore)
	binApp := app_service_binary.NewBinaryAppService(binStore)
	textApp := app_service_text.NewTextAppService(textStore)
	cardApp := app_service_card.NewCardAppService(cardStore)

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
