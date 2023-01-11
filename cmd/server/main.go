package main

import (
	"GophKeeper/internal/server/app_services/app_service_auth"
	"GophKeeper/internal/server/app_services/app_service_binary"
	"GophKeeper/internal/server/app_services/app_service_card"
	"GophKeeper/internal/server/app_services/app_service_credential"
	"GophKeeper/internal/server/app_services/app_service_text"
	"GophKeeper/internal/server/server_grpc/services/grpc_service_auth"
	"GophKeeper/internal/server/server_grpc/services/grpc_service_binary"
	"GophKeeper/internal/server/server_grpc/services/grpc_service_card"
	"GophKeeper/internal/server/server_grpc/services/grpc_service_cred"
	"GophKeeper/internal/server/server_grpc/services/grpc_service_text"
	"GophKeeper/internal/storage/binary_store"
	"GophKeeper/internal/storage/card_store"
	"GophKeeper/internal/storage/credential_store"
	"GophKeeper/internal/storage/text_store"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"

	"go.uber.org/zap"

	"GophKeeper/internal/server"
	"GophKeeper/internal/server/server_grpc"
	"GophKeeper/internal/server/server_grpc/interceptors"
	"GophKeeper/internal/storage/auth_store"
	"GophKeeper/pkg/logzap"
)

var (
	_ = (*server_grpc.ServerGRPC)(nil)
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
		textStore = text_store.NewPostgresStorage(db)
		binStore = binary_store.NewPostgresStorage(db)

		credStore = credential_store.NewMemoryStorage()
		cardStore = card_store.NewMemoryStorage()

	} else {
		authStore = auth_store.NewMemoryStorage()
		credStore = credential_store.NewMemoryStorage()
		binStore = binary_store.NewMemoryStorage()
		textStore = text_store.NewMemoryStorage()
		cardStore = card_store.NewMemoryStorage()
	}

	// Создание сервисов приложения
	authApp := app_service_auth.NewAuthService(authStore, app_service_auth.WithSecretKey(cfg.SecretKey))
	credApp := app_service_credential.NewCredentialAppService(credStore)
	binApp := app_service_binary.NewBinaryAppService(binStore)
	textApp := app_service_text.NewTextAppService(textStore)
	cardApp := app_service_card.NewCardAppService(cardStore)

	// Создание gRPC сервисов
	authRPC := grpc_service_auth.NewAuthServiceRPC(authApp)
	credRPC := grpc_service_cred.NewCredServiceRPC(credApp)
	binRPC := grpc_service_binary.NewBinaryServiceRPC(binApp)
	textRPC := grpc_service_text.NewTextServiceRPC(textApp)
	cardRPC := grpc_service_card.NewCardServiceRPC(cardApp)

	validate := interceptors.NewValidateInterceptor(cfg.SecretKey)

	// Создание сервера
	grpcServer, err := server_grpc.NewServer(
		cfg.AddrGRPC,
		validate,
		server_grpc.WithAuthServiceRPC(authRPC),
		server_grpc.WithCredServiceRPC(credRPC),
		server_grpc.WithBinaryServiceRPC(binRPC),
		server_grpc.WithTextServiceRPC(textRPC),
		server_grpc.WithCardServiceRPC(cardRPC),
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
