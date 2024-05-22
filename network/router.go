package network

// HandleFunc 路由处理函数
type HandleFunc func()

// group 路由分组
type group struct {
	prefix     string
	handlerMap map[string]HandleFunc
}

// router 路由
type router struct {
	groups map[string]*group
}
