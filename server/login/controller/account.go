package controller

import (
	"github.com/zerodot618/zerokk-sg-server/network"
	"github.com/zerodot618/zerokk-sg-server/server/login/proto"
)

var DefaultAccount = &Account{}

type Account struct{}

func (a *Account) Router(r *network.Router) {
	g := r.Group("account")

	g.AddRouter("login", a.login)
}

func (a *Account) login(req *network.WsMsgReq, rsp *network.WsMsgRsp) {
	rsp.Body.Code = 0
	loginRes := &proto.LoginRsp{
		Username: "admin",
		Password: "admin",
		Session:  "session",
		UId:      1,
	}
	rsp.Body.Msg = loginRes
}
