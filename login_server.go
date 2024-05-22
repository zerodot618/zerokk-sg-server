package main

import (
	"github.com/zerodot618/zerokk-sg-server/config"
	"github.com/zerodot618/zerokk-sg-server/network"
)

func main() {
	host := config.File.MustValue("login_server", "host", "127.0.0.1")
	port := config.File.MustValue("login_server", "port", "8004")

	server := network.NewServer(host + ":" + port)
	server.Start()
}
