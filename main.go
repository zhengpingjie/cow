package main

import (
	"net/http"
	_ "net/http/pprof"
	"robot/app/config"
	"robot/app/grpc/grpcserver"
	"robot/app/route"
)
func main(){
//	engine := gin.Default()
	// 设置路由
	//runtime.GOMAXPROCS(runtime.NumCPU())
	  //http路由
	//route.HttpRouter(engine)
	  //注册websocket路由
	route.WsInit()
	//执行定时清理任务
	route.Task()
	go route.WsRoute()
	//grpc路由
	go grpcserver.Init()
	//engine.Run(config.ApiPort)
	http.ListenAndServe(config.HttpPort, nil)
}