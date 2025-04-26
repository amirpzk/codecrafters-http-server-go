package main

import (
	"github.com/codecrafters-io/http-server-starter-go/app/server"
)

func main() {
	config := server.NewConfig()
	srv := server.NewServer(config)
	srv.Start()
}
