package main

import (
	"GophKeeper/internal/storage/auth_store"
	"GophKeeper/internal/storage/data_store/binary_store"
	"GophKeeper/internal/storage/data_store/card_store"
	"GophKeeper/internal/storage/data_store/credential_store"
	"GophKeeper/internal/storage/data_store/text_store"
	"database/sql"
	"fmt"
	"github.com/golang-migrate/migrate/v4"
	"os"
	"os/signal"
	"syscall"

	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"

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

	dsn := "user=postgres password=postgres dbname=GophKeeper sslmode=disable"
	db := PostgresDB(dsn)
	if err := migrateFor(db.DB, "postgres"); err != nil && err != migrate.ErrNoChange {
		zap.L().Error("error migrate", zap.Error(err))
	} else {
		zap.L().Info("Success apply migrations")
	}

	cfg := newConfig()

	validate := interceptors.NewValidateInterceptor(cfg.SecretKey)

	// Создание хранилищ
	authStore := auth_store.NewMemoryStorage()
	credStore := credential_store.NewMemoryStorage()
	binStore := binary_store.NewMemoryStorage()
	textStore := text_store.NewMemoryStorage()
	cardStore := card_store.NewMemoryStorage()

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

func PostgresDB(dsn string) *sqlx.DB {

	logger := zap.L()

	db, errOpen := sqlx.Open("postgres", dsn)
	if errOpen != nil {
		logger.Error("failed to connect to the database: %s\n", zap.Error(errOpen))
	}

	if err := db.Ping(); err != nil {
		logger.Error("connection to DB created, but Ping returned error: %s\n", zap.Error(err))
	}

	logger.Info("Success connect to database")

	return db
}

func migrateFor(db *sql.DB, driverDB string) error {
	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		return err
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://db/migrations",
		driverDB, driver)

	if err != nil {
		return err
	}

	return m.Up()
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
