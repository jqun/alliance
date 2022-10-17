package handler

import (
	"alliance/common/net"
	"alliance/proto-message/pb"
	"log"
)

func ScAccountLogin(data interface{}, c net.Conner) {
	resp, ok := data.(*pb.ScAccountLogin)
	if !ok {
		return
	}
	if resp.Result != 1 {
		log.Printf("account[%v] login exception, logout please, then login new account", resp.RoleName)
		return
	}
	c.SetRoleName(resp.RoleName)
	log.Printf("account[%v] login server success", resp.RoleName)
}

func ScAccountLogout(data interface{}, c net.Conner) {
	_, ok := data.(*pb.ScAccountLogout)
	if !ok {
		return
	}
	roleName := c.GetRoleName()
	c.SetRoleName("")
	log.Printf("account[%v] logout success",  roleName)
}