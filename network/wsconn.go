package network

// ReqBody 请求
type ReqBody struct {
	Seq   int64       `json:"seq"`
	Name  string      `json:"name"`
	Msg   interface{} `json:"msg"`
	Proxy string      `json:"proxy"`
}

// RspBody 响应
type RspBody struct {
	Seq  int64       `json:"seq"`
	Name string      `json:"name"`
	Code int         `json:"code"`
	Msg  interface{} `json:"msg"`
}

// WsMsgReq ws请求
type WsMsgReq struct {
	Body *ReqBody
	Conn WSConn
}

// WsMsgRsp ws响应
type WsMsgRsp struct {
	Body *RspBody
}

// WSConn websocket连接
type WSConn interface {
	SetProperty(key string, value interface{})
	GetProperty(key string) (interface{}, error)
	RemoveProperty(key string)
	Addr() string
	Push(name string, data interface{})
}
