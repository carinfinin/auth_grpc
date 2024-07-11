package main

import (
	"auth/internal/config"
	"fmt"
)

func main() {
	cfg := config.MustLoad()
	fmt.Println(cfg)

	// logger
	// app init
	// grpc server start
}
