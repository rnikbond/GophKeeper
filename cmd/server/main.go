package main

import (
	"GophKeeper/internal/server"
	"fmt"
	"log"
	"time"
)

func main() {

	cfg := server.NewConfig()
	if err := cfg.ParseArgs(); err != nil {
		log.Fatalln(err.Error())
	}

	fmt.Printf("Server started on %s at: %s\n", cfg.AddrGRPC, time.Now().Format("02-01-2006 15:04:05"))
	time.Sleep(1 * time.Second)
	fmt.Printf("Server stopped at: %s\n", time.Now().Format("02-01-2006 15:04:05"))

}
