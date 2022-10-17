package main

import (
	"alliance/client/_client"
	"alliance/client/config"
	"alliance/client/handler"
	"runtime"
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	handler.InitHandler()
	_client.Client = _client.NewClient(config.GetClientConfig().ServerUrl)
	_client.Client.Run()
}
