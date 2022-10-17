package handler

import (
	"alliance/common/net"
	"alliance/proto-message/pb"
	pbId "alliance/proto-message/pb-id"
	"alliance/server/module/role"
	"alliance/server/module/server"
	"log"
)

func CsAccountLogin(data interface{}, c net.Conner) {
	req := data.(*pb.CsAccountLogin)
	user, ok := server.SeverInfo.GetUser(req.RoleName)
	routeMsg := &pb.ScAccountLogin{Result: 1, RoleName: req.RoleName}
	if c.GetRoleName() != "" { // 重复登录
		routeMsg.Result = 2
		net.SendMsg(pbId.ScAccountLoginId, routeMsg, c)
		return
	}
	if !ok { //不存在 直接创角
		user = role.NewUser(req.RoleName)
		server.SeverInfo.StoreUser(user)
		log.Printf("create account[%v] user success", req.RoleName)
	} else {
		log.Printf("account[%v] already exist and user login success", req.RoleName)
	}
	c.SetRoleName(req.RoleName)
	server.SeverInfo.NetSever.StoreConner(user.Role.Name, c)
	net.SendMsg(pbId.ScAccountLoginId, routeMsg, c)
}

func CsAccountLogout(data interface{}, c net.Conner) {
	_, ok := data.(*pb.CsAccountLogout)
	if !ok {
		return
	}
	_, ok = server.SeverInfo.GetUser(c.GetRoleName())
	if !ok { //不存在
		return
	}
	log.Printf("account[%v] logout server", c.GetRoleName())
	//server.SeverInfo.UserLogout(c.GetRoleName())
	c.SetRoleName("")
	net.SendMsg(pbId.ScAccountLogoutId, &pb.ScAccountLogout{}, c)
}
