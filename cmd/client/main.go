package main

import (
	"GophKeeper/internal/client/clientgrpc"
	"GophKeeper/internal/server"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
)

var (
	_ = (*clientgrpc.ClientGRPC)(nil)
)

var (
	buildVersion = "N/A"
	buildDate    = "N/A"
	buildCommit  = "N/A"
)

func main() {

	cfg := newConfig()
	cli := newClient(cfg)

	if err := cli.Connect(); err != nil {
		log.Fatalf("server connection error: %v\n", err)
	}

	if err := cli.Login(); err != nil {

		if err = cli.Register(); err != nil {
			log.Fatalf("error say hello: %v\n", err)
		}

		fmt.Println("Success Register")
	}

	fmt.Println("Success Login")

	done := make(chan os.Signal)
	signal.Notify(done, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)
	<-done

	if err := cli.Disconnect(); err != nil {
		log.Printf("connection error from the server: %v\n ", err)
	}

	fmt.Println("Goodbye...")
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

// newClient Создание объекта клиента
func newClient(cfg *server.Config) *clientgrpc.ClientGRPC {

	return clientgrpc.NewClient(cfg.AddrGRPC)
}
