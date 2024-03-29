package main

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"

	"github.com/fatih/color"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"GophKeeper/internal/client"
	"GophKeeper/internal/client/app_services/app_service_auth"
	"GophKeeper/internal/client/app_services/app_service_binary"
	"GophKeeper/internal/client/app_services/app_service_card"
	"GophKeeper/internal/client/app_services/app_service_cred"
	"GophKeeper/internal/client/app_services/app_service_text"
	"GophKeeper/internal/client/grpc_services/grpc_service_auth"
	"GophKeeper/internal/client/grpc_services/grpc_service_binary"
	"GophKeeper/internal/client/grpc_services/grpc_service_card"
	"GophKeeper/internal/client/grpc_services/grpc_service_cred"
	"GophKeeper/internal/client/grpc_services/grpc_service_text"
	"GophKeeper/pkg/logzap"
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

	conn, err := grpc.Dial(cfg.AddrGRPC, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		logger.Fatal("failed gRPC connect", zap.Error(err))
	}

	cli := newClient(conn, cfg)
	cli.Start()

	if err := conn.Close(); err != nil {
		logger.Fatal("failed gRPC disconnect", zap.Error(err))
	}
}

func init() {

	fmt.Printf("Build version: %s\n", buildVersion)
	fmt.Printf("Build date: %s\n", buildDate)
	fmt.Printf("Build commit: %s\n", buildCommit)
}

func newConfig() *client.Config {
	cfg := client.NewConfig()
	if err := cfg.ParseArgs(); err != nil {
		logger := zap.L()
		logger.Fatal("failed run gRPC server: %v\n", zap.Error(err))
	}

	return cfg
}

func newClient(conn *grpc.ClientConn, cfg *client.Config) *client.Client {

	pubKey := publicKey(cfg.PublicKey)
	privKey := privateKey(cfg.PrivateKey)

	if pubKey != nil && privKey != nil {
		color.Green("Encoding data: enabled")
	} else {
		color.Yellow("Encoding data: disabled")
	}

	rpcAuth := grpc_service_auth.NewService(conn)
	rpcText := grpc_service_text.NewService(conn)
	rpcBin := grpc_service_binary.NewService(conn)
	rpcCred := grpc_service_cred.NewService(conn)
	rpcCard := grpc_service_card.NewService(conn)

	authApp := app_service_auth.NewService(rpcAuth, app_service_auth.WithSalt(cfg.Salt))
	textApp := app_service_text.NewService(rpcText, app_service_text.WithPublicKey(pubKey), app_service_text.WithPrivateKey(privKey))
	binApp := app_service_binary.NewService(rpcBin, app_service_binary.WithPublicKey(pubKey), app_service_binary.WithPrivateKey(privKey))
	credApp := app_service_cred.NewService(rpcCred, app_service_cred.WithPublicKey(pubKey), app_service_cred.WithPrivateKey(privKey))
	cardApp := app_service_card.NewService(rpcCard, app_service_card.WithPublicKey(pubKey), app_service_card.WithPrivateKey(privKey))

	return client.NewClient(authApp,
		client.WithService(textApp),
		client.WithService(binApp),
		client.WithService(credApp),
		client.WithService(cardApp))
}

func publicKey(key []byte) *rsa.PublicKey {

	block, _ := pem.Decode(key)
	if block == nil {
		return nil
	}

	pubKey, errKey := x509.ParsePKIXPublicKey(block.Bytes)
	if errKey != nil {
		return nil
	}

	switch pub := pubKey.(type) {
	case *rsa.PublicKey:
		return pub
	default:
		return nil
	}
}

func privateKey(key []byte) *rsa.PrivateKey {

	block, _ := pem.Decode(key)
	if block == nil {
		return nil
	}

	privKey, errParse := x509.ParsePKCS1PrivateKey(block.Bytes)
	if errParse != nil {
		return nil
	}

	return privKey
}
