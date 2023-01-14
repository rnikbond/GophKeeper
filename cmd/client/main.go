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
	clientGrpc "GophKeeper/internal/client/client_grpc"
	"GophKeeper/internal/client/client_grpc/services/auth_service"
	"GophKeeper/internal/client/client_grpc/services/binary_service"
	"GophKeeper/internal/client/client_grpc/services/card_service"
	"GophKeeper/internal/client/client_grpc/services/cred_service"
	"GophKeeper/internal/client/client_grpc/services/text_service"
	"GophKeeper/pkg/logzap"
)

var (
	_ = (*clientGrpc.ClientGRPC)(nil)
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

func newClient(conn *grpc.ClientConn, cfg *client.Config) *clientGrpc.ClientGRPC {

	pubKey := publicKey(cfg.PublicKey)
	privKey := privateKey(cfg.PrivateKey)

	if pubKey != nil && privKey != nil {
		color.Green("Encoding data: enabled")
	} else {
		color.Yellow("Encoding data: disabled")
	}

	authServ := auth_service.NewService(conn, auth_service.WithSalt(cfg.Salt))
	textServ := text_service.NewService(conn, text_service.WithPublicKey(pubKey), text_service.WithPrivateKey(privKey))
	binServ := binary_service.NewService(conn, binary_service.WithPublicKey(pubKey), binary_service.WithPrivateKey(privKey))
	credServ := cred_service.NewService(conn, cred_service.WithPublicKey(pubKey), cred_service.WithPrivateKey(privKey))
	cardServ := card_service.NewService(conn, card_service.WithPublicKey(pubKey), card_service.WithPrivateKey(privKey))

	return clientGrpc.NewClient(
		authServ,
		clientGrpc.WithService(textServ),
		clientGrpc.WithService(binServ),
		clientGrpc.WithService(credServ),
		clientGrpc.WithService(cardServ))
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
