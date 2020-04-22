package main

import (
	log "github.com/sirupsen/logrus"
	"myth/go-essential/app"
	"myth/go-essential/base/rpc/server"
	"myth/go-essential/net/rpc/warden"
	pb "myth/go-example/proto"
	"myth/go-example/server/common"
	"myth/go-example/server/handler"
	"myth/go-example/server/manager"
)

func main() {
	p := app.GetMythApp()
	p.Config = &common.Config{}
	p.Run(
		app.WithLogger(),
		app.With(func(mpp *app.MythApp) error {
			log.Info("With")
			return nil
		}),
		app.WithManager(func(mpp *app.MythApp) app.Manager {
			log.Info("WithManager")
			manager := manager.NewManager()
			return manager
		}),
		app.WithCronTab(func(mpp *app.MythApp) error {
			log.Info("WithCronTab")
			return nil
		}),
		app.WithCronTab(func(mpp *app.MythApp) error {
			log.Info("WithCronTab")
			return nil
		}),
		//app.WithRpcClient(func(client client.Client, mpp *app.MythApp) error {
		//	log.Info("WithRpcClient")
		//	return nil
		//}),
		app.WithRpcServer(func(srv server.Server, mpp *app.MythApp) error {
			log.Info("WithRpcServer")
			server := srv.(*warden.Server)
			conf := mpp.Config.(*common.Config)

			//中间件测试
			//server.UseUnary(handler.ExampleAuthFunc())
			server.RegisterRpc()

			hdr := handler.NewHandler(conf)
			pb.RegisterGreeterServer(server.RpcServer(), hdr)
			return nil
		}),
	)
}
