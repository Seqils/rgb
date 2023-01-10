package main

import (
	"rgb/internal/server"
	"rgb/internal/conf"
	"rgb/internal/cli"
)

func main() {
	cli.Parse()
	server.Start(conf.NewConfig())
}