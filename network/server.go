package network

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

// server web服务
type server struct {
	addr   string  // 服务地址
	router *router // 路由
}

// NewServer 创建服务器
func NewServer(addr string) *server {
	return &server{
		addr: addr,
	}
}

// 启动服务
func (s *server) Start() {
	// 定义路由
	http.HandleFunc("/", s.wsHandler)

	// 监听服务
	err := http.ListenAndServe(s.addr, nil)
	if err != nil {
		panic(err)
	}
}

// http 升级 websocket协议的配置
var wsUpgrader = websocket.Upgrader{
	// 允许所有 CORS 跨域请求
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func (s *server) wsHandler(w http.ResponseWriter, r *http.Request) {
	// 1. 将 http 协议设计为 websocket 协议
	wsConn, err := wsUpgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Fatal("websocket 服务连接出错:", err)
	}
	log.Println("websocket 服务连接成功")
	// 2. 创建 websocket 服务
	// websocket 通道建立后，客户端服务端都可以收发消息
	// 发消息时，把消息当作路由来处理，消息是有格式的
	// 需要先定义消息的格式
	wsServer := NewWsServer(wsConn)
	wsServer.Router(s.router)
	// 3. 启动websocket服务
	wsServer.Start()
}
