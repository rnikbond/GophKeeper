package main

import (
	"GophKeeper/internal/server"
	"GophKeeper/internal/server/servergrpc"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
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

	cfg := newConfig()
	serv := newServer(cfg)

	serv.Start()

	done := make(chan os.Signal)
	signal.Notify(done, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)
	<-done

	serv.Stop()
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
		log.Fatalf("failed run: %v\n", err.Error())
	}

	return cfg
}

// newServer Создание объекта сервера
func newServer(cfg *server.Config) *servergrpc.ServerGRPC {

	serv, err := servergrpc.NewServer(cfg.AddrGRPC)
	if err != nil {
		log.Fatalf("failed run: %v\n", err.Error())
	}

	return serv
}
