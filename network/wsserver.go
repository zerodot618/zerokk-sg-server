package network

import (
	"encoding/json"
	"fmt"
	"log"
	"sync"

	"github.com/forgoer/openssl"
	"github.com/gorilla/websocket"
	"github.com/zerodot618/zerokk-sg-server/utils"
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
		// 1. data 解压
		data, err = utils.UnZip(data)
		if err != nil {
			fmt.Println("解压数据出错: ", err)
			continue
		}
		// 2. 解密消息
		secretKey, err := w.GetProperty("secretKey")
		if err == nil {
			// 有加密
			key := secretKey.(string)
			// 解密
			d, err := utils.AesCBCDecrypt(data, []byte(key), []byte(key), openssl.ZEROS_PADDING)
			if err != nil {
				log.Println("解密数据出错: ", err)
				// 出错 发起握手
				// w.Handshake()
			} else {
				data = d
			}
		}
		// 3. data 转为 body
		body := &ReqBody{}
		err = json.Unmarshal(data, body)
		if err != nil {
			log.Println("数据格式有误: ", err)
		} else {
			// 获取到数据
			req := &WsMsgReq{
				Body: body,
				Conn: w,
			}
			rsp := &WsMsgRsp{
				Body: &RspBody{
					Seq:  req.Body.Seq,
					Name: body.Name,
				},
			}
			w.router.Run(req, rsp)
			w.outChan <- rsp
		}

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
