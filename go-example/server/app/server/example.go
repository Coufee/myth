package main

import (
	log "github.com/sirupsen/logrus"
	"myth/go-essential/app"
	"myth/go-essential/base/rpc/client"
	"myth/go-essential/base/rpc/server"
	"myth/go-essential/net/rpc/warden"
	pb "myth/go-example/proto"
	"myth/go-example/server/handler"
	"myth/go-example/server/manager"
)

func main() {
	p := app.GetMythApp()
	p.Run(
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
		app.WithRpcClient(func(client client.Client, mpp *app.MythApp) error {
			log.Info("WithRpcClient")
			return nil
		}),
		app.WithRpcServer(func(srv server.Server, mpp *app.MythApp) error {
			log.Info("WithRpcServer")
			server := srv.(*warden.Server)
			hdr := handler.NewHandler()
			pb.RegisterHelloServiceServer(server.GetServer(), hdr)
			return nil
		}),
	)
}
