package main

import (
	"apiProxy/config"
	"apiProxy/router"
	"giga"
)

func main() {
	r := giga.NewEngine()
	router.InitRpcClient()
	router.InitRouter(r)

	r.Run(config.DefaultConfig.App.Name, config.DefaultConfig.App.Addr)
}
