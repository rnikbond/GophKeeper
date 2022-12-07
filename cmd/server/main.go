package main

import (
	"fmt"
	"time"
)

func main() {

	fmt.Printf("Server started at: %s\n", time.Now().Format("02-01-2006 15:04:05"))
	time.Sleep(5 * time.Second)
	fmt.Printf("Server stopped at: %s\n", time.Now().Format("02-01-2006 15:04:05"))

}
