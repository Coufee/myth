package main

import (
	"myth/go-essential/app"
	"myth/go-essential/base/rpc/client"
	log "myth/go-essential/log/logc"
	"myth/go-essential/net/rpc/warden"
	"myth/go-example/middle/common"
	"myth/go-example/middle/handler/client_handler"
	"myth/go-example/middle/manager"
)

func main() {
	p := app.GetMythApp()
	p.Config = &common.Config{}
	p.Run(
		app.WithManager(func(mpp *app.MythApp) app.Manager {
			log.Info("WithManager")
			manager := manager.NewManager()
			return manager
		}),
		app.WithCronTab(func(mpp *app.MythApp) error {
			log.Info("WithCronTab")
			return nil
		}),
		app.WithRpcClient(func(client client.Client, mpp *app.MythApp) error {
			log.Info("WithRpcClient")
			c := client.(*warden.Client)

			//中间件
			i:=0
			c.UseUnary(client_handler.ExampleAuthFunc(&i))
			hdr := client_handler.NewHandler(c)

			//for ;;i++{
				hdr.SayHello("hello")
			//}

			return nil
		}),
		//app.WithRpcServer(func(srv server.Server, mpp *app.MythApp) error {
		//	log.Info("WithRpcServer")
		//	server := srv.(*warden.Server)
		//
		//	//中间件测试
		//	server.RegisterRpc()
		//
		//	hdr := server_handler.NewHandler()
		//	pb.RegisterGreeterServer(server.RpcServer(), hdr)
		//	return nil
		//}),
	)
}
