package network

// HandleFunc 路由处理函数
type HandleFunc func()

// group 路由分组
type group struct {
	prefix     string                // 路由前缀
	handlerMap map[string]HandleFunc // 路由处理函数
}

// router 路由
type router struct {
	groups map[string]*group // 路由分组
}

func (r *router) Run(req *WsMsgReq, rsp *WsMsgRsp) {
	
}
