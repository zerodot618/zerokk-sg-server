package main

import (
	"github.com/zerodot618/zerokk-sg-server/config"
	"github.com/zerodot618/zerokk-sg-server/network"
)

func main() {
	// 读取配置文件 conf.ini 中的内容
	host := config.File.MustValue("login_server", "host", "127.0.0.1")
	port := config.File.MustValue("login_server", "port", "8004")

	// 启动服务
	server := network.NewServer(host + ":" + port)
	server.Start()
}
