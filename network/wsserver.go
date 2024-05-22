package network

import (
	"fmt"
	"log"
	"sync"

	"github.com/gorilla/websocket"
)

// wsServer websocket 服务
type wsServer struct {
	wsConn       *websocket.Conn        // 连接
	router       *router                // 路由
	outChan      chan *WsMsgRsp         // 写队列
	Seq          int64                  // 序号
	property     map[string]interface{} // 属性
	propertyLock sync.RWMutex           // 属性锁
}

func NewWsServer(wsConn *websocket.Conn) *wsServer {
	return &wsServer{
		wsConn:   wsConn,
		outChan:  make(chan *WsMsgRsp, 1000),
		property: make(map[string]interface{}),
		Seq:      0,
	}
}

// Start 通道建立，收发消息需要一直监听
func (w *wsServer) Start() {
	// 启动读写数据的处理逻辑
	go w.readMsgLoop()
	go w.writeMsgLoop()
}

func (w *wsServer) readMsgLoop() {
	defer func() {
		if err := recover(); err != nil {
			log.Fatal(err)
			w.Close()
		}
	}()
	// 先读到客服端发送过来的数据,然后进行处理,然后再写回给客服端
	for {
		_, data, err := w.wsConn.ReadMessage()
		if err != nil {
			fmt.Println("读取消息出错: ", err)
			break
		}
		// 将收到的数据进行路由
		fmt.Println(data)
	}
	w.Close()
}

func (w *wsServer) Close() {
	_ = w.wsConn.Close()
}

func (w *wsServer) writeMsgLoop() {
	for {
		select {
		case msg := <-w.outChan:
			fmt.Println("write msg: ", msg)
		}
	}
}

func (w *wsServer) Router(router *router) {
	w.router = router
}

func (w *wsServer) SetProperty(key string, value interface{}) {
	w.propertyLock.Lock()
	defer w.propertyLock.Unlock()
	w.property[key] = value
}

func (w *wsServer) GetProperty(key string) (interface{}, error) {
	w.propertyLock.RLock()
	defer w.propertyLock.RUnlock()
	return w.property[key], nil
}

func (w *wsServer) RemoveProperty(key string) {
	w.propertyLock.Lock()
	defer w.propertyLock.Unlock()
	delete(w.property, key)
}

func (w *wsServer) Addr() string {
	return w.wsConn.RemoteAddr().String()
}

func (w *wsServer) Push(name string, data interface{}) {
	rsp := &WsMsgRsp{
		Body: &RspBody{
			Seq:  0,
			Name: name,
			Msg:  data,
		},
	}
	w.outChan <- rsp
}
