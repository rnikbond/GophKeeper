package main

import (
	"GophKeeper/internal/server"
	"GophKeeper/internal/server/grpc_server"
	"log"
	"os"
	"os/signal"
	"syscall"
)

var (
	_ = (*grpc_server.GRPCServer)(nil)
)

func main() {

	cfg := newConfig()
	serv := newServer(cfg)
	done := make(chan os.Signal)

	serv.Start()

	signal.Notify(done, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)
	<-done

	serv.Stop()
}

func newServer(cfg *server.Config) *grpc_server.GRPCServer {

	serv, err := grpc_server.NewServer(cfg.AddrGRPC)
	if err != nil {
		log.Fatalf("failed run: %v\n", err.Error())
	}

	return serv
}

func newConfig() *server.Config {
	cfg := server.NewConfig()
	if err := cfg.ParseArgs(); err != nil {
		log.Fatalf("failed run: %v\n", err.Error())
	}

	return cfg
}
