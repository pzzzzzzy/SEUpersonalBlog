package main

import (
	"personal-blog/config"
	"personal-blog/internal/server"
)

func main() {
	cfg := config.LoadConfig()
	srv := server.NewServer(cfg)
	srv.Start()
}