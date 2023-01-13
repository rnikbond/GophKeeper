package main

import (
	clientGrpc "GophKeeper/internal/client/client_grpc"
	"GophKeeper/internal/client/client_grpc/services/auth_service"
	"GophKeeper/internal/client/client_grpc/services/text_service"
	"GophKeeper/internal/server"
	"GophKeeper/pkg/logzap"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

//
//var (
//	_ = (*clientgrpc.ClientGRPC)(nil)
//)
//
//var (
//	buildVersion = "N/A"
//	buildDate    = "N/A"
//	buildCommit  = "N/A"
//)

func main() {

	logzap.ConfigZapLogger()

	logger := zap.L()
	cfg := newConfig()

	conn, err := grpc.Dial(cfg.AddrGRPC, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		logger.Fatal("failed gRPC connect", zap.Error(err))
	}

	cli := newClient(conn)
	cli.Start()

	if err := conn.Close(); err != nil {
		logger.Fatal("failed gRPC disconnect", zap.Error(err))
	}

}

//func init() {
//
//	fmt.Printf("Build version: %s\n", buildVersion)
//	fmt.Printf("Build date: %s\n", buildDate)
//	fmt.Printf("Build commit: %s\n", buildCommit)
//}
//

func newConfig() *server.Config {
	cfg := server.NewConfig()
	if err := cfg.ParseArgs(); err != nil {
		logger := zap.L()
		logger.Fatal("failed run gRPC server: %v\n", zap.Error(err))
	}

	return cfg
}

func newClient(conn *grpc.ClientConn) *clientGrpc.ClientGRPC {

	authServ := auth_service.NewService(conn)
	textServ := text_service.NewService(conn)

	return clientGrpc.NewClient(authServ, clientGrpc.WithTextService(textServ))
}
