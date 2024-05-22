package login

import (
	"github.com/zerodot618/zerokk-sg-server/network"
	"github.com/zerodot618/zerokk-sg-server/server/login/controller"
)

var Router = network.NewRouter()

func Init() {
	initRouter()
}

func initRouter() {
	controller.DefaultAccount.Router(Router)
}
