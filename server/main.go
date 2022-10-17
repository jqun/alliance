package main

import (
	network "alliance/common/net"
	"alliance/server/config"
	"alliance/server/handler"
	"alliance/server/module/server"
	"runtime"
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	handler.InitHandler()
	s := network.NewServer(config.GetServerConfig().App.Addr)
	server.SeverInfo.SetNetServer(s)
	server.SeverInfo.Run()
}
