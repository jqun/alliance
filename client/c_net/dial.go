package c_net

import (
	network "alliance/common/net"
	"log"
	"net"
	"time"
)

func NewClientConner(addr string) network.Conner {
	conn := netDial(addr)
	conner := network.NewConner(conn)
	return conner
}

func netDial(addr string) net.Conn {
	connectInterval := time.Second
	for {
		conn, err := net.Dial("tcp", addr)
		if conn != nil {
			return conn
		}
		if err != nil {
			log.Printf("dial addr[%v] server err,error is [%v]", addr, err.Error())
		}
		connectInterval *= 3
		time.Sleep(connectInterval)
	}
}
