package main

import (
	log "github.com/sirupsen/logrus"
	"maybe/article/lottery/common"

	"myth/go-essential/app"
	"myth/go-essential/base/rpc/server"
	"myth/go-essential/net/rpc/warden"
	"myth/go-example/manager"
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
		//app.WithHttpServer(func(e *gin.Engine, mpp *app.MythApp) error {
		//	log.Info("WithHttpServer")
		//	return nil
		//}),
		app.WithRpcClient(func(mpp *app.MythApp) error {
			log.Info("WithRpcClient")
			return nil
		}),
		app.WithRpcServer(func(srv server.Server, mpp *app.MythApp) error {
			log.Info("WithRpcServer")
			server := srv.(*warden.Server)
			server.Start()

			return nil
		}),
	)
}
