package network

import "strings"

// HandleFunc 路由处理函数
type HandleFunc func(req *WsMsgReq, rsp *WsMsgRsp)

// group 路由分组
type group struct {
	prefix     string                // 路由前缀
	handlerMap map[string]HandleFunc // 路由处理函数
}

func (g *group) exec(name string, req *WsMsgReq, rsp *WsMsgRsp) {
	if h, ok := g.handlerMap[name]; ok {
		h(req, rsp)
	}
}

func (g *group) AddRouter(name string, handleFunc HandleFunc) {
	g.handlerMap[name] = handleFunc
}

// Router 路由
type Router struct {
	group []*group // 路由分组
}

func NewRouter() *Router {
	return &Router{}
}

func (r *Router) Group(prefix string) *group {
	g := &group{
		prefix:     prefix,
		handlerMap: make(map[string]HandleFunc),
	}
	r.group = append(r.group, g)
	return g
}

func (r *Router) Run(req *WsMsgReq, rsp *WsMsgRsp) {
	// req.Body.Name 路径
	strs := strings.Split(req.Body.Name, ".")
	prefix := ""
	name := ""
	if len(strs) == 2 {
		prefix = strs[0]
		name = strs[1]
	}

	for _, g := range r.group {
		if g.prefix == prefix {
			g.exec(name, req, rsp)
		}
	}
}
